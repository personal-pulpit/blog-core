package model

import (
	"time"
)


type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `gorm:"size:25;NOT NULL"`
	LastName  string `gorm:"size:25;NOT NULL"`
	Email     string `gorm:"unique;size:50;NOT NULL"`
	Biography string `gorm:"type:text;size:500;NOT NULL"`
}

func NewUser(firtsname, lastname, email, biography string) *User {
	return &User{
		FirstName: firtsname,
		LastName:  lastname,
		Email:     email,
		Biography: biography,
	}
}
