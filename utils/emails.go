package utils

import (
	"fmt"
	"os"

	m "github.com/keighl/mandrill"
)

var (
	emailToken    string
	userEmailFrom string
	userNameFrom  string
)

func InitEmail() {
	emailToken = os.Getenv("EMAIL_TOKEN")
	userEmailFrom = os.Getenv("USER_EMAIL_FROM")
	userNameFrom = os.Getenv("USER_NAME_FROM")
}

func SendEmail(fromEmail, fromName, toEmail, toName, subject, html, text string) error {
	client := m.ClientWithKey(emailToken)

	message := &m.Message{}
	message.AddRecipient(toEmail, toName, "to")
	message.FromEmail = fromEmail
	message.FromName = fromName
	message.Subject = subject
	message.HTML = html
	message.Text = text

	_, err := client.MessagesSend(message)
	return err
}

func SendChangePasswordEmail(email, name, token string) error {
	html :=
		`<h1>You forgot the password<h1>
    <p>this is the token to create a new one %s, this token is valid for the next hour</p>`
	text :=
		`You forgot the password\n
  this is the token to create a new one %s, this token is valid for the next hour`
	subject := "Change Password"
	return SendEmail(userEmailFrom, userNameFrom, email, name, subject, fmt.Sprintf(html, token), fmt.Sprintf(text, token))
}
