package model

type Role int

const (
	UserRole Role = iota + 1
	AdminRole
)

type Auth struct {
	ID                  uint
	HashedPassword      string `gorm:"size:300;NOT NULL"`
	FailedLoginAttempts int    `gorm:"failed_login_attempts;NOT NULL"`
	AccountLockedUntil  int64  `gorm:"account_locked_until;NOT NULL"`
	EmailVerified       bool   `gorm:"email_verified;NOT NULL"`
	Role                Role   `gorm:"default:1;NOT NULL"`
}

func NewAuth(ID uint, password string, role Role) *Auth {
	return &Auth{
		ID:                  ID,
		HashedPassword:      password,
		FailedLoginAttempts: 0,
		AccountLockedUntil:  0,
		EmailVerified:       false,
		Role:                role,
	}
}
