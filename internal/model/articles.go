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
	AuthorId  ID     `gorm:"NOT NULL"`
}

func NewArticle(title, content string, authorID ID) *Article {
	return &Article{
		Title:    title,
		Content:  content,
		AuthorId: authorID,
	}
}
