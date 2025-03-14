package model

import (
	"time"
)


type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `gorm:"size:25;NOT NULL"`
	LastName  string `gorm:"size:25;NOT NULL"`
	Email     string `gorm:"unique;size:50;NOT NULL"`
	Biography string `gorm:"type:text;size:500;NOT NULL"`
	Articles  []Article `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	Comments []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func NewUser(firstName, lastName, email, biography string) *User {
	return &User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Biography: biography,
	}
}
