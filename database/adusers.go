package database

import (
  "strings"
	"log"
	"o365-attack-toolkit/model"
	_ "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func GetADUsers() []model.ADUser{

	var users []model.ADUser


	rows, err := db.Query(model.GetADUsersQuery)
	user := model.ADUser{}

	if err != nil{
		log.Println("Error : " + err.Error())
	}
	for rows.Next() {
    var phones string
		err := rows.Scan(&user.ID,&phones,&user.DisplayName,&user.GivenName,&user.Mail,&user.MobilePhone,&user.PreferredLanguage,&user.Surname,&user.UserPrincipalName)
    user.BusinessPhones = strings.Split(phones, ",")
		if err != nil {
			log.Fatal(err)
		}
		users = append(users,user)
	}

	return users
}


func InsertADUser(user model.ADUser){

	tx, _ := db.Begin()
	stmt, err_stmt := tx.Prepare(model.InsertADUserQuery)

	if err_stmt != nil {
		log.Fatal(err_stmt)
	}
	_, err := stmt.Exec(user.ID,strings.Join(user.BusinessPhones,","),user.DisplayName,user.GivenName,user.Mail,user.MobilePhone,user.PreferredLanguage,user.Surname,user.UserPrincipalName)
	tx.Commit()
	if err != nil{
		log.Printf("ERROR: %s",err)
	}

}
