package models

type User struct {
	Id uint `gorm:"primaryKey"`
	Username string `gorm:"type:string;unique;size:15;not null"`
	Password string `gorm:"type:string;size:100;not null"`
}
