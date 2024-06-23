package model

import "time"
type ID string
type Base struct {
	ID        ID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
