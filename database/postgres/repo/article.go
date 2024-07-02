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

func (a *articlePostgresRepo) GetAll() ([]*model.Article, error) {
	var articles = []*model.Article{}

	if err := a.postgresCLI.Find(articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}
func (a *articlePostgresRepo) GetArticle(filters map[string]interface{}) (*model.Article, error) {
	var article = &model.Article{}

	query := a.postgresCLI

	for key, value := range filters {
		query = query.Where(key, value)
	}

	if err := query.Find(article).Error; err != nil {
		return nil, err
	}

	return article, nil
}
func (a *articlePostgresRepo) GetArticleByTitle(title string) (*model.Article, error) {
	filter := map[string]interface{}{
		"title": title,
	}

	article, err := a.GetArticle(filter)

	if err != nil {
		return nil, err
	}

	return article, nil
}
func (a *articlePostgresRepo) GetArticleById(id model.ID) (*model.Article, error) {
	filter := map[string]interface{}{
		"id": id,
	}

	article, err := a.GetArticle(filter)

	if err != nil {
		return nil, err
	}

	return article, nil
}
func (a *articlePostgresRepo) Create(articleModel *model.Article) (*model.Article, error) {
	err := a.postgresCLI.Create(&articleModel).Error

	if err != nil {
		return nil, err
	}

	return articleModel, nil
}
func (a *articlePostgresRepo) UpdateByID(ID model.ID, title, content string) (*model.Article, error) {
	var article *model.Article

	err := a.postgresCLI.First(&article, ID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return article, ErrArticleNotFound
		}

		return article, err
	}

	article.Title = title
	article.Content = content

	err = a.postgresCLI.Save(&article).Error

	if err != nil {
		return article, err
	}

	return article, err
}
func (a *articlePostgresRepo) DeleteByID(ID model.ID) error {
	var article *model.Article

	err := a.postgresCLI.Delete(article, ID).Error

	if err != nil {
		return err
	}

	return nil
}
