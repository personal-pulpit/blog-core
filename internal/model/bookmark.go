package model

import "time"

type Bookmark struct {
	ID        uint     `gorm:"primaryKey"`
	UserID    uint     `gorm:"foreignKey:UserID"`
	ArticleID uint     `gorm:"foreignKey:ArticleID"`
	Article   *Article `gorm:"foreignKey:ArticleID"`
	CreatedAt time.Time
}

func NewBookmark(userID, articleID uint) *Bookmark {
	return &Bookmark{
		UserID:    userID,
		ArticleID: articleID,
	}
}
