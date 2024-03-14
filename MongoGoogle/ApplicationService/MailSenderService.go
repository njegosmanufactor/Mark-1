package ApplicationService

import (
	model "MongoGoogle/Model"
	conn "MongoGoogle/Repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Sends an email for account verification to the provided email address.
func SendMail(email string) {
	// Sender data.
	from := "nemanja.ranitovic@manufactoryteam.io"
	var password, pass_err = os.LookupEnv("GOOGLE_MAIL_PASSWORD")
	if !pass_err {
		log.Fatal("Google_mail_password not declared in .env file!")
	}
	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("Controller/pages/MailTemplate.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Account verification \n%s\n\n", mimeHeaders)))
	link := "http://localhost:3000/verify/{email}"
	link = strings.Replace(link, "{email}", email, 1)
	t.Execute(&body, struct {
		Message string
	}{
		Message: link,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")

}

// Sends an email for ownership transfer to the provided email address and sends a response if the user is not found.
func SendOwnershipMail(email string, res http.ResponseWriter) {

	collection := conn.GetClient().Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	var result model.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return
		}
	} else {
		// Sender data.
		from := "nemanja.ranitovic@manufactoryteam.io"
		var password, pass_err = os.LookupEnv("GOOGLE_MAIL_PASSWORD")
		if !pass_err {
			log.Fatal("Google_mail_password not declared in .env file!")
		}
		// Receiver email address.
		to := []string{
			email,
		}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		// Authentication.
		auth := smtp.PlainAuth("", from, password, smtpHost)

		t, _ := template.ParseFiles("Controller/pages/OwnershipMailTemplate.html")

		var body bytes.Buffer

		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		body.Write([]byte(fmt.Sprintf("Subject: Ownership transfer \n%s\n\n", mimeHeaders)))
		link := "http://localhost:3000/transferOwnership/feedback/{email}"
		link = strings.Replace(link, "{email}", email, 1)
		t.Execute(&body, struct {
			Message string
		}{
			Message: link,
		})

		// Sending email.
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Email Sent!")
	}
}

// Sends an email invitation to the provided email address for joining a company identified by the given company ID.
func SendInvitationMail(id string, email string) {
	// Sender data.
	from := "nemanja.ranitovic@manufactoryteam.io"
	var password, pass_err = os.LookupEnv("GOOGLE_MAIL_PASSWORD")
	if !pass_err {
		log.Fatal("Google_mail_password not declared in .env file!")
	}
	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("Controller/pages/InviteTemplate.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Company invitation \n%s\n\n", mimeHeaders)))
	link := "http://localhost:3000/inviteConfirmation/{id}"
	link = strings.Replace(link, "{id}", id, 1)
	t.Execute(&body, struct {
		Message string
	}{
		Message: link,
	})

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!")
}
