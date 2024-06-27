package email

import (
	"blog/pkg/auth_manager"
	"blog/pkg/email_manager"
	"blog/utils/random"
	"fmt"
	"time"
)
const EmailVerficationCodeExp = time.Minute * 2
type EmailService interface {
	SendVerficationEmailCode(recipientEmail string) error
}
type emailOpts struct {
	authManager  auth_manager.AuthManager
	emailManager email_manager.EmailManager
	uniqueID string
}

func NewEmailService(authManager auth_manager.AuthManager, emailManager email_manager.EmailManager) EmailService {
	return &emailOpts{
		authManager:  authManager,
		emailManager: emailManager,
	}
}
func (s *emailOpts) SendVerficationEmailCode(recipientEmail string) error {
	s.uniqueID = fmt.Sprintf("%d",random.GenerateUniqueID())
	otp,err := s.authManager.SetOTP(s.uniqueID,EmailVerficationCodeExp)
	if err != nil{
		return ErrSendingVerificationEmailFaild
	}
	err = s.emailManager.SendVerificationEmail(recipientEmail,otp)
	if err != nil{
		return ErrSendingVerificationEmailFaild
	}
	return nil
}
func (s *emailOpts)CompareCodeAndOTP(code string)error{
	otp,err := s.authManager.GetOTP(s.uniqueID)
	if err != nil{
		return 	ErrCheckingCodeFaild
	}
	if code != otp{
		return ErrInvalidCode
	}
	return nil
}
