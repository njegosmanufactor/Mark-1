package Service

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

var smtpHost = "smtp.gmail.com"
var smtpPort = "587"

type EmailParams struct {
	Subject string
	Body    string
}

func sendEmail(from string, to []string, auth smtp.Auth, params EmailParams) {
	message := fmt.Sprintf("Subject: %s\n%s", params.Subject, params.Body)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
		return
	}
	log.Println("Email Sent!")
}

func generateLink(baseLink, placeholder, value string) string {
	return strings.Replace(baseLink, placeholder, value, 1)
}

func SetMailSender(email string) (string, []string, smtp.Auth) {
	from := "nemanja.ranitovic@manufactoryteam.io"
	password, passErr := os.LookupEnv("GOOGLE_MAIL_PASSWORD")
	if !passErr {
		log.Fatal("Google_mail_password not declared in .env file!")
	}
	to := []string{email}
	auth := smtp.PlainAuth("", from, password, smtpHost)
	return from, to, auth
}

func LoadHTMLTemplate(filePath string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func SendEmailWithHTMLTemplate(email, subject, filePath string, data interface{}) {
	from, to, auth := SetMailSender(email)

	tmpl, err := LoadHTMLTemplate(filePath)
	if err != nil {
		fmt.Println("Error loading HTML template:", err)
		return
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

	// Set MIME headers to indicate HTML content.
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Concatenate MIME headers and HTML body.
	fullBody := mimeHeaders + body.String()

	sendEmail(from, to, auth, EmailParams{Subject: subject, Body: fullBody})
}

func SendMail(email string) {
	link := generateLink("http://localhost:3000/auth/verify/{email}", "{email}", email)
	SendEmailWithHTMLTemplate(email, "Account verification", "Controller/pages/MailTemplate.html", struct{ Message string }{Message: link})
}

func SendOwnershipMail(transferID, email string, res http.ResponseWriter) {
	link := generateLink("http://localhost:3000/user/transferOwnership/feedback/{transferID}", "{transferID}", transferID)
	SendEmailWithHTMLTemplate(email, "Ownership transfer", "Controller/pages/OwnershipMailTemplate.html", struct{ Message string }{Message: link})
}

func SendInvitationMail(id, email string) {
	link := generateLink("http://localhost:3000/user/inviteConfirmation/{id}", "{id}", id)
	SendEmailWithHTMLTemplate(email, "Company invitation", "Controller/pages/InviteTemplate.html", struct{ Message string }{Message: link})
}

func SendPasswordChangeLink(id, email string) {
	link := generateLink("http://localhost:3000/user/forgotPassword/callback/PAGETOREDIRECTTO/{id}", "{id}", id)
	SendEmailWithHTMLTemplate(email, "Change password request", "Controller/pages/ChangePasswordTemplate.html", struct{ Message string }{Message: link})
}

func SendMagicLink(email string) {
	link := generateLink("http://localhost:3000/auth/confirmMagicLink/{email}", "{email}", email)
	SendEmailWithHTMLTemplate(email, "Application confirmation", "Controller/pages/MagicLink.html", struct{ Message string }{Message: link})
}

func SendPasswordLessCode(email, code string) {
	SendEmailWithHTMLTemplate(email, "Application confirmation", "Controller/pages/PasswordLess.html", struct{ Message string }{Message: code})
}
