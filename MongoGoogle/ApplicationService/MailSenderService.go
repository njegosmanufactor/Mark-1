package ApplicationService

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"
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
func SendOwnershipMail(transferId string, email string, res http.ResponseWriter) {
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
	link := "http://localhost:3000/transferOwnership/feedback/{transferId}"
	//ovde treba da ide id transakcije
	link = strings.Replace(link, "{transferId}", transferId, 1)
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

func SendPasswordChangeLink(id string, email string) {
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

	t, _ := template.ParseFiles("Controller/pages/ChangePasswordTemplate.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Change password request \n%s\n\n", mimeHeaders)))
	link := "http://localhost:3000/forgotPassword/callback/PAGETOREDIRECTTO/{id}"
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

func SendMagicLink(email string) {
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

	t, _ := template.ParseFiles("Controller/pages/MagicLink.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Application confirmation \n%s\n\n", mimeHeaders)))
	link := "http://localhost:3000/magicLink/{email}"
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
