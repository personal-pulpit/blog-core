package model

import (
	"blog/utils/random"
	"fmt"
	"time"
)

type Role int

const (
	UserRole Role = iota + 1
	AdminRole
)

type User struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string `gorm:"size:25;NOT NULL"`
	LastName  string `gorm:"size:25;NOT NULL"`
	Email     string `gorm:"unique;size:50;NOT NULL"`
	Biography string `gorm:"type:text;size:500;NOT NULL"`
	Role      Role   `gorm:"default:1;NOT NULL"`
}

func NewUser(firtsname, lastname, email, biography string, role Role) *User {
	id := fmt.Sprintf("%d", random.GenerateId())

	return &User{
		ID:        id,
		FirstName: firtsname,
		LastName:  lastname,
		Email:     email,
		Biography: biography,
		Role:      role,
	}
}
