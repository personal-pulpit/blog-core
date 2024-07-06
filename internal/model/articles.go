package model

import (
	"blog/utils/random"
	"fmt"
	"time"
)

type Article struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	Title     string `gorm:"size:100;NOT NULL"`
	Content   string `gorm:"text;NOT NULL"`
	AuthorId  ID     `gorm:"NOT NULL"`
}

func NewArticle(title, content string, authorID ID) *Article {
	id := fmt.Sprintf("%d", random.GenerateId())
	return &Article{
		ID:       id,
		Title:    title,
		Content:  content,
		AuthorId: authorID,
	}
}
