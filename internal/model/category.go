package model

import "time"


type Category struct{
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"size:100;NOT NULL"`
	Articles  []Article `gorm:"many2many:article_categories;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
}

func NewCategory(name string)*Category{
	return &Category{
		Name: name,
	}
}