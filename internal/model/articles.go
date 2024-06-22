package model

type Article struct {
	Base
	Title    string `gorm:"size:100;NOT NULL"`
	Content  string `gorm:"text;NOT NULL"`
	AuthorId uint   `gorm:"NOT NULL"`
}
