package server

import (
	"fmt"
  "strings"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"o365-attack-toolkit/database"
	"o365-attack-toolkit/model"
	"os"
	"path/filepath"
)

// This will contain the template functions

func ExecuteTemplate(w http.ResponseWriter,page model.Page, templatePath string)  {
  tpl := template.New("")
  tpl.Funcs(template.FuncMap{"StringsJoin": strings.Join})
  _,err := tpl.ParseFiles("templates/main.html",templatePath)
	if err != nil{
		log.Fatal(err)
	}
  tpl.ExecuteTemplate(w,"layout",page)
}

func ExecuteSingleTemplate(w http.ResponseWriter,page model.Page, templatePath string){
	tpl,err := template.ParseFiles(templatePath)
	if err != nil{
		log.Println(err)
	}

	err = tpl.Execute(w,page)
	if err != nil{
		log.Println(err)
	}

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Users"
	Page.UserList = database.GetUsers();
	ExecuteTemplate(w,Page,"templates/users.html")
}

func GetADUsers(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "AD Users"
	Page.ADUserList = database.GetADUsers();
	ExecuteTemplate(w,Page,"templates/adusers.html")
}

func GetUserEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Emails"

	vars := mux.Vars(r)
	id := vars["id"]

	Page.Email = id
	Page.EmailList = database.GetEmailsByUser(id)
	ExecuteTemplate(w,Page,"templates/emails.html")
}

func GetAllEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Emails"
	Page.Email = "all"
	Page.EmailList = database.GetAllEmails()
	ExecuteTemplate(w,Page,"templates/emails.html")
}


func SearchUserEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Emails"
	vars := mux.Vars(r)
	user := vars["id"]
	searchKey := r.FormValue("search")
	Page.EmailList = database.SearchUserEmails(user,searchKey)

	Page.Email = user
	ExecuteTemplate(w,Page,"templates/emails.html")
	}



func SearchEmails(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}

	searchKey := r.FormValue("search")
	Page.Title = "Search result for " + searchKey

	Page.EmailList = database.SearchEmails(searchKey)
	ExecuteTemplate(w,Page,"templates/emails.html")
}


func GetUserFiles(w http.ResponseWriter, r *http.Request) {
	Page := model.Page{}
	Page.Title = "Files"

	vars := mux.Vars(r)
	email := vars["email"]

	folderDir := fmt.Sprintf("./downloads/%s",email)
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
	ExecuteSingleTemplate(w,Page,"templates/files.html")
	//ExecuteTemplate(w,Page,"templates/files.html")
}

func GetUserFile(){}


func GetAbout(w http.ResponseWriter, r *http.Request){
	Page := model.Page{}
	Page.Title = "About"
	ExecuteTemplate(w,Page,"templates/about.html")
}


func GetEmail(w http.ResponseWriter, r *http.Request){

	var Page model.Page
	vars := mux.Vars(r)
	email := vars["id"]
	Page.Mail = database.GetEmail(email)
	ExecuteSingleTemplate(w,Page,"templates/email.html")
}
