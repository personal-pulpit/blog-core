package mysql_repository

import (
	"blog/internal/model"
	"blog/internal/repository"

	"errors"

	"gorm.io/gorm"
)

type articleMysqlRepo struct {
	mysqlClient *gorm.DB
}

var (
	ErrArticleNotFound = errors.New("article not found")
)

func NewArticleMysqlRepo(mysqlCLI *gorm.DB) repository.ArticleMysqlRepository {
	return &articleMysqlRepo{
		mysqlClient: mysqlCLI,
	}
}

func (a *articleMysqlRepo) Create(authorID model.ID, title, content string) (model.Article, error) {
	var article model.Article
	article.Title = title
	article.Content = content
	article.AuthorId = authorID
	tx := NewTx(a.mysqlClient)
	err := tx.Create(&article).Error
	if err != nil {
		return article, err
	}
	tx.Commit()
	return article, nil
}
func (a *articleMysqlRepo) UpdateByID(ID, title, content string) (model.Article, error) {
	var article model.Article
	err := a.mysqlClient.First(&article, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return article, ErrArticleNotFound
		}
		return article, err
	}
	article.Title = title
	article.Content = content
	tx := NewTx(a.mysqlClient)
	err = tx.Save(&article).Error
	if err != nil {
		return article, err
	}
	tx.Commit()
	return article, err
}
func (a *articleMysqlRepo) DeleteByID(ID string) error {
	var article model.Article
	tx := NewTx(a.mysqlClient)
	err := tx.Delete(&article, ID).Error
	if err != nil {
		return err
	}
	return nil
}
