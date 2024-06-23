package model

type Article struct {
	Base
	Title    string `gorm:"size:100;NOT NULL"`
	Content  string `gorm:"text;NOT NULL"`
	AuthorId ID   `gorm:"NOT NULL"`
}

func  NewArticle(title,content string,authorID ID)*Article{
	return &Article{
		Title: title,
		Content: content,
		AuthorId: authorID,
	}
}