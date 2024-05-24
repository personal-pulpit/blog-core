package models
type Role uint
const(
	UserRole Role = 1
	AdminRole Role = 2

)
type User struct {
	Id uint `gorm:"primaryKey"`
	Fristname string `gorm:"size:25;not null"`
	Lastname string `gorm:"size:25;not null"`
	Username string `gorm:"unique;size:25;not null"`
	Password string `gorm:"size:150;not null"`
	Email string `gorm:"unique;size:50;not null"`
	PhoneNumber string `gorm:"unique;size:11;not null"`
	Role uint `gorm:"default:1;not null"`
}
