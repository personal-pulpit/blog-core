package model

type Auth struct {
	ID             ID `gorm:"type:string;NOT NULL"`
	HashedPassword string `gorm:"size:300;NOT NULL"`
	FailedLoginAttempts int    `gorm:"failed_login_attempts;NOT NULL"`
	AccountLockedUntil  int64  `gorm:"account_locked_until;NOT NULL"`
	EmailVerified       bool   `gorm:"email_verified;NOT NULL"`
}
func NewAuth(ID ID,hashedPassword string)*Auth{
	return &Auth{
		ID: ID,
		HashedPassword: hashedPassword,
		FailedLoginAttempts: 0,
		AccountLockedUntil: 0,
		EmailVerified: false,
	}
}