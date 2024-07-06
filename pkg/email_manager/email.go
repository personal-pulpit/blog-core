package email

import (
	"blog/config"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"path/filepath"
	"runtime"
)

type EmailService interface {
	sendEmail(msg *emailMessage) error
	readTemplate(template string) *template.Template
	SendWelcomeEmail(recipientEmail, name string) error
	SendResetPasswordEmail(recipientEmail, url, name, exp string) error
	SendVerificationEmail(recipientEmail, otp string) error
}

func templatesDirPath() string {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Error in generating env dir")
	}

	return filepath.Dir(f)
}

type emailOtp struct {
	Configs *config.Email
}

type emailMessage struct {
	tmplateFileName string
	receiver        []string
	args            interface{}
	subject         string
}

func NewEmailService(Configs *config.Email) EmailService {
	return &emailOtp{Configs}
}

func newEmailMessage(tmplFileName, subject string,args interface{}, reciver []string) *emailMessage {
	return &emailMessage{
		tmplateFileName: tmplFileName,
		subject:         subject,
		args:            args,
		receiver:        reciver,
	}
}

func (s *emailOtp) readTemplate(templFileName string) *template.Template {
	templFileNameFullAddres := fmt.Sprintf("%s/%s/%s",templatesDirPath(),"templates",templFileName)

	tpl := template.Must(template.ParseFiles(templFileNameFullAddres))

	return tpl
}

func (s *emailOtp) sendEmail(msg *emailMessage) error {
	if config.GetEnv() == config.Development {
		log.Println("Email sent")
		log.Println(msg)
		return nil
	}

	template := s.readTemplate(msg.tmplateFileName)

	var body bytes.Buffer
	err := template.Execute(&body, msg.args)
	
	if err != nil {
		return err
	}

	emailMessage := fmt.Sprintf("Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+body.String(), msg.subject)
	
	auth := smtp.PlainAuth("", s.Configs.SenderEmail, s.Configs.Password, s.Configs.Host)
	err = smtp.SendMail(s.Configs.Host+":"+s.Configs.Port, auth, s.Configs.SenderEmail, msg.receiver, []byte(emailMessage))
	
	if err != nil {
		return err
	}
	
	log.Println("Verification code email sent successfully!")
	
	return nil
}

func (s *emailOtp) SendWelcomeEmail(recipientEmail, name string) error {
	msg := newEmailMessage(
		"welcome.html",
		"Welcome",
		struct {Name string}{Name: name},
		[]string{recipientEmail},
	)
	
	err := s.sendEmail(msg)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *emailOtp) SendVerificationEmail(recipientEmail, otp string) error {
	msg := newEmailMessage(
		"verification_code.html",
		"Verify Account",
		struct {OTP string}{OTP: otp},
		[]string{recipientEmail},
	)
	
	err := s.sendEmail(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *emailOtp) SendResetPasswordEmail(recipientEmail, url, name, exp string) error {
	msg := newEmailMessage(
		"submit_reset_password.html",
		"Reset Password",
		struct {Name,RecipientEmail,URL,EXP string}{Name: name,RecipientEmail: recipientEmail,URL: url,EXP: exp},
		[]string{recipientEmail},
	)

	err := s.sendEmail(msg)
	if err != nil {
		return err
	}

	return nil
}
