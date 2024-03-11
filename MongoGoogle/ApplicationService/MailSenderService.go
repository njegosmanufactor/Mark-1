package ApplicationService

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
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
