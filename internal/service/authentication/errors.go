package authentication

import "errors"

var (
	ErrInvalidEmailOrPassword   = errors.New("invalid email or password")
	ErrAccessDenied             = errors.New("access denied")
	ErrEmailNotVerified         = errors.New("email not verified")
	ErrAccountLocked            = errors.New("account locked")
	ErrCodeIsInvalid            = errors.New("code is invalid")
	ErrCheckPasswordHash        = errors.New("failed to check password hash")
	ErrDecodeAccessToken        = errors.New("failed to decode access token")
	ErrDecodeRefreshToken       = errors.New("failed to decode refresh token")
	ErrHashPassword             = errors.New("failed to hash password")
	ErrVerifyEmail              = errors.New("failed to verify email")
	ErrRefreshToken             = errors.New("failed to refresh token")
	ErrGenerateToken            = errors.New("failed to generate token")
	ErrSendVerificationCode     = errors.New("failed to send verification code")
	ErrSendResetPassword        = errors.New("failed to send reset password")
	ErrSendWelcome              = errors.New("failed to send welcome")
	ErrGenerateVerificationCode = errors.New("failed to generate verification code")
	ErrDestroyRefreshToken      = errors.New("failed to destroy refresh token")
	ErrSetTokenIntoBlackList    = errors.New("failed to set token into black list")
	ErrChangePassword           = errors.New("failed to change password")
	ErrLockAccount              = errors.New("failed to lock account")
	ErrGenerateUniqueID         = errors.New("failed to generate unique ID")
)
