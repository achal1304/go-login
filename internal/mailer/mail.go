package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
	"time"

	"github.com/achal1304/go-login/internal/data"
	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New() Mailer {
	dialer := mail.NewDialer("smtp.mailtrap.io", 25, "0a268b42179455", "7d5bebbf92ebae")
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: "achal20000@gmail.com",
	}
}

func (m Mailer) Send(recipient, templateFile string, userId int, email string) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	token, err := data.GenerateToken(userId, 15*time.Minute, email)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resetLink := "http://localhost:3000/resetPassword/" + token.Plaintext

	data := map[string]interface{}{
		"resetLink": resetLink,
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", "Reset Password Link")
	msg.SetBody("text/plain", plainBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
