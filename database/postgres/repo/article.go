package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"

	"errors"

	"gorm.io/gorm"
)

type articlePostgresRepo struct {
	postgresCLI *gorm.DB
}

var (
	ErrArticleNotFound = errors.New("article not found")
)

func NewArticlePostgresRepo(postgresCLI *gorm.DB) repository.ArticlePostgresRepository {
	return &articlePostgresRepo{
		postgresCLI: postgresCLI,
	}
}

func (a *articlePostgresRepo) Create(authorID model.ID, title, content string) (model.Article, error) {
	var article model.Article
	article.Title = title
	article.Content = content
	article.AuthorId = authorID
	tx := NewTx(a.postgresCLI)
	err := tx.Create(&article).Error
	if err != nil {
		return article, err
	}
	tx.Commit()
	return article, nil
}
func (a *articlePostgresRepo) UpdateByID(ID, title, content string) (model.Article, error) {
	var article model.Article
	err := a.postgresCLI.First(&article, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return article, ErrArticleNotFound
		}
		return article, err
	}
	article.Title = title
	article.Content = content
	tx := NewTx(a.postgresCLI)
	err = tx.Save(&article).Error
	if err != nil {
		return article, err
	}
	tx.Commit()
	return article, err
}
func (a *articlePostgresRepo) DeleteByID(ID string) error {
	var article model.Article
	tx := NewTx(a.postgresCLI)
	err := tx.Delete(&article, ID).Error
	if err != nil {
		return err
	}
	return nil
}
