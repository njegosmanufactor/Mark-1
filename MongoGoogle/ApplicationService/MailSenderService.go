package ApplicationService

import (
	"crypto/tls"
	"fmt"
	"strings"

	gomail "gopkg.in/mail.v2"
)

func SendMail(email string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "nemanja.ranitovic@manufactoryteam.io")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Account verification")
	link := "http://localhost:3000/verify/{email}"
	link = strings.Replace(link, "{email}", email, 1)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", link)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "nemanja.ranitovic@manufactoryteam.io", "cwcn trol loos svbr")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
