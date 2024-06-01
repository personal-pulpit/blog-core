package models

type Role uint

const (
	UserRole  Role = 1
	AdminRole Role = 2
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
	Role        uint   `gorm:"default:1;NOT NULL"`
}
