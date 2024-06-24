package model

import "time"

type (
	ID  =  string
	Base struct {
		ID        ID `gorm:"primaryKey"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
