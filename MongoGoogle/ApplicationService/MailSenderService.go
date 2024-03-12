package ApplicationService

import (
	conn "MongoGoogle/MongoDB"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SendMail(email string) {
	// Sender data.
	from := "nemanja.ranitovic@manufactoryteam.io"
	password := "cwcn trol loos svbr"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("LoginRegister/pages/MailTemplate.html")

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

func SendOwnershipMail(email string, res http.ResponseWriter) {
	collection := conn.Client.Database("UserDatabase").Collection("Users")
	filter := bson.M{"Email": email}
	var result conn.ApplicationUser
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(res).Encode("Didnt find user!")
			return
		}
	} else {
		// Sender data.
		from := "nemanja.ranitovic@manufactoryteam.io"
		password := "cwcn trol loos svbr"

		// Receiver email address.
		to := []string{
			email,
		}

		// smtp server configuration.
		smtpHost := "smtp.gmail.com"
		smtpPort := "587"

		// Authentication.
		auth := smtp.PlainAuth("", from, password, smtpHost)

		t, _ := template.ParseFiles("LoginRegister/pages/OwnershipMailTemplate.html")

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

func SendInvitationMail(email string, compnayID string) {
	// Sender data.
	from := "nemanja.ranitovic@manufactoryteam.io"
	password := "cwcn trol loos svbr"

	// Receiver email address.
	to := []string{
		email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("LoginRegister/pages/InviteTemplate.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Company invitation \n%s\n\n", mimeHeaders)))
	link := "http://localhost:3000/inviteConfirmation/{companyID}/{email}"
	link = strings.Replace(link, "{email}", email, 1)
	link = strings.Replace(link, "{companyID}", compnayID, 1)
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
