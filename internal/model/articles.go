package model

import (
		"time"
)

type Article struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `gorm:"size:100;NOT NULL"`
	Content   string `gorm:"text;NOT NULL"`
	AuthorID  uint
	Author *User  `gorm:"foreignKey:AuthorID"`
}

func NewArticle(title, content string, authorID uint) *Article {
	return &Article{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
	}
}
