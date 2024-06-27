package email_manager

import (
	"blog/config"
	"bytes"
	"fmt"
	"html/template"
	"log"

	"net/smtp"
)
type EmailManager interface {
	SendVerificationEmail(recipientEmail string, otp string) error
}
var verificationTemplate *template.Template

func init() {
	var err error
	verificationTemplate, err = template.ParseFiles("template/verification_code.html")
	if err != nil {
		panic(err)
	}
}

type EmailOpts struct {
	configs *config.Email
}

func NewEmailService(configs *config.Email) EmailManager {
	return &EmailOpts{configs}
}
func (e *EmailOpts) SendVerificationEmail(recipientEmail string, otp string) error {
	currentENV := config.GetEnv()
	if currentENV == config.Production {
		var buffer bytes.Buffer
		err := verificationTemplate.Execute(&buffer, struct{ Code string }{Code: otp})
		if err != nil {
			return err
		}
		body := buffer.String()
		message := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Verification Code\r\n\r\n%s", "Blog", recipientEmail, body))

		auth := smtp.PlainAuth("", e.configs.SenderEmail, e.configs.Password, e.configs.Host)
		err = smtp.SendMail(e.configs.Host+":"+e.configs.Port, auth, e.configs.SenderEmail, []string{recipientEmail}, message)
		if err != nil {
			return err
		}
		log.Println("Verification code email sent successfully!")
		return nil
	}
	log.Println("Email sent")
	return nil
}
