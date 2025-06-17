package model

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"foreignKey:UserID"`
	ArticleID uint   `gorm:"foreignKey:ArticleID"`
	Content   string `gorm:"NOT NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewComment(userID, articleID uint, content string) *Comment {
	return &Comment{
		UserID:    userID,
		ArticleID: articleID,
		Content:   content,
	}
}
