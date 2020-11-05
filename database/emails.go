package database

import (
	_ "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"o365-attack-toolkit/model"
)

func GetEmailsByUser(email string) []model.Mail{

	var mails []model.Mail



	rows, err := db.Query(model.GetUserMailsQuery,email)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&mail.Id,&mail.User,&mail.Subject,&mail.SenderEmail,&mail.SenderName,&mail.HasAttachments,&mail.BodyPreview,&mail.BodyType,&mail.BodyContent)
		if err != nil {
			log.Fatal(err)
		}
		mails = append(mails,mail)
	}

	return mails
}

func GetAllEmails() []model.Mail{

	var mails []model.Mail



	rows, err := db.Query(model.GetMailsQuery)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&mail.Id,&mail.User,&mail.Subject,&mail.SenderEmail,&mail.SenderName,&mail.HasAttachments,&mail.BodyPreview,&mail.BodyType,&mail.BodyContent)
		if err != nil {
			log.Fatal(err)
		}
		mails = append(mails,mail)
	}

	return mails


}

func InsertEmail(mail model.Mail){

	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.InsertMailQuery)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(mail.Id,mail.User,mail.Subject,mail.SenderEmail,mail.SenderName,mail.HasAttachments,mail.BodyPreview,mail.BodyType,mail.BodyContent)
	tx.Commit()
	if err != nil{
		log.Printf("ERROR: %s",err)
	}

}

func SearchUserEmails(email string,searchKey string) []model.Mail {
	var mails []model.Mail

	searchKey = "%" + searchKey + "%"

	rows, err := db.Query(model.SearchUserMailsQuery,email,searchKey)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&mail.Id,&mail.User,&mail.Subject,&mail.SenderEmail,&mail.SenderName,&mail.HasAttachments,&mail.BodyPreview,&mail.BodyType,&mail.BodyContent)
		if err != nil {
			log.Fatal(err)
		}
		mails = append(mails,mail)
	}

	return mails
}


func SearchEmails(searchKey string) []model.Mail {
	var mails []model.Mail

	searchKey = "%" + searchKey + "%"

	rows, err := db.Query(model.SearchEmailQuery,searchKey)
	mail := model.Mail{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&mail.Id,&mail.User,&mail.Subject,&mail.SenderEmail,&mail.SenderName,&mail.HasAttachments,&mail.BodyPreview,&mail.BodyType,&mail.BodyContent)
		if err != nil {
			log.Fatal(err)
		}
		mails = append(mails,mail)
	}

	return mails
}


func GetEmail(id string) model.Mail {

	row := db.QueryRow(model.GetEmailQuery,id)
	mail := model.Mail{}

	err := row.Scan(&mail.Id,&mail.User,&mail.Subject,&mail.SenderEmail,&mail.SenderName,&mail.HasAttachments,&mail.BodyPreview,&mail.BodyType,&mail.BodyContent)

	if err != nil {
		//It's empty
		log.Println(err)
		return mail
	}

	return mail
}
