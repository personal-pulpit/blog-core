package model


type Role int

const (
	UserRole Role = iota +1
	AdminRole 
)

type User struct {
	Base
	Firstname   string `gorm:"size:25;NOT NULL"`
	Lastname    string `gorm:"size:25;NOT NULL"`
	Username    string `gorm:"unique;size:20;NOT NULL"`
	Password    string `gorm:"size:150;NOT NULL"`
	Email       string `gorm:"unique;size:50;NOT NULL"`
	PhoneNumber string `gorm:"unique;size:11;NOT NULL"`
	Biography   string `gorm:"type:text;size:500;NOT NULL"`
	Role        Role   `gorm:"default:1;NOT NULL"`
}
func NewUser(firtsname,lastname,username,email,phoneNumber,biography string,role Role)*User{
	return &User{
		Firstname: firtsname,
		Lastname: lastname,
		Username: username,
		Email: email,
		PhoneNumber: phoneNumber,
		Biography: biography,
		Role: role,
	}
}