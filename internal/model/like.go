package model

import "time"

type Like struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"foreignKey:UserID"`
	ArticleID uint `gorm:"foreignKey:ArticleID"`
	Article *Article
	CreatedAt time.Time
}

func NewLike(userID, articleID uint) *Like {
	return &Like{
		UserID:    userID,
		ArticleID: articleID,
	}
}