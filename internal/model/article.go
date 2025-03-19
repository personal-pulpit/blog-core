package model

import (
	"time"
)

type Article struct {
    ID        uint   `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Title     string `gorm:"size:100;NOT NULL"`
    Content   string `gorm:"type:text;NOT NULL"`
    AuthorID  uint
    Author    *User   `gorm:"foreignKey:AuthorID"`
    Comments  []Comment `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;"`
	Categories []Category `gorm:"many2many:article_categories;"`
}

type ArticleFilter struct {
	Title       *string 
	PublishedAt *time.Time 
	PublishedAtGT *time.Time
	PublishedAtLT *time.Time
}

func NewArticle(title, content string, authorID uint,categories []Category) *Article {
	return &Article{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Categories: categories,
	}
}
