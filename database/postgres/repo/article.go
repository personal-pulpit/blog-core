package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"errors"

	"gorm.io/gorm"
)

type articlePostgresRepo struct {
	postgresCLI *gorm.DB
}

func NewArticlePostgresRepo(postgresCLI *gorm.DB) repository.ArticleRepository {
	return &articlePostgresRepo{
		postgresCLI: postgresCLI,
	}
}

func (a *articlePostgresRepo) GetAll(ctx context.Context) ([]model.Article, error) {
	var articles = []model.Article{}

	if err := a.postgresCLI.WithContext(ctx).Preload("Author").Preload("Comments").Preload("Categories").Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("get all articles: %w: %v", repository.ErrDatabase,err)
	}

	return articles, nil
}

func (a *articlePostgresRepo) SearchArticles(ctx context.Context, filter *model.ArticleFilter) ([]model.Article, error) {
	var articles []model.Article
	query := a.postgresCLI.WithContext(ctx).Model(&model.Article{})

	var conditions []string
	var args []interface{}

	if filter.Title != nil {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+*filter.Title+"%")
	}

	if filter.PublishedAt != nil {
		conditions = append(conditions, "created_at = ?")
		args = append(args, *filter.PublishedAt)
	}

	if filter.PublishedAtGT != nil {
		conditions = append(conditions, "created_at > ?")
		args = append(args, *filter.PublishedAtGT)
	}

	if filter.PublishedAtLT != nil {
		conditions = append(conditions, "created_at < ?")
		args = append(args, *filter.PublishedAtLT)
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	err := query.Preload("Author").Preload("Comments").Preload("Categories").Find(&articles).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrArticleNotFound
		}
		return nil, fmt.Errorf("search articles: args:  %v\n%w: %v",args, repository.ErrDatabase,err )
	}

	return articles, nil
}

func (a *articlePostgresRepo) GetArticleByID(ctx context.Context, ID uint) (*model.Article, error) {
	article := new(model.Article)

	err := a.postgresCLI.WithContext(ctx).Preload("Author").Preload("Comments").Preload("Categories").First(article, ID).Error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrArticleNotFound
		}
		return nil, fmt.Errorf("get article by id: article ID:  %d\n%w: %v",ID, repository.ErrDatabase,err )
	}

	return article, nil
}

func (a *articlePostgresRepo) Create(ctx context.Context, articleModel *model.Article) (*model.Article, error) {
	err := a.postgresCLI.WithContext(ctx).Create(&articleModel).Error
	if err != nil {
		return nil, fmt.Errorf("create article: %v \n%w: %v",articleModel, repository.ErrDatabase, err)
	}

	return articleModel, nil
}

func (a *articlePostgresRepo) UpdateByID(ctx context.Context, ID uint, title, content string) (*model.Article, error) {
	article := new(model.Article)

	article,err  := a.GetArticleByID(ctx,ID)
	if err != nil {
		return nil,err
	}

	article.Title = title
	article.Content = content

	err = a.postgresCLI.WithContext(ctx).Save(article).Error
	if err != nil {
		return nil, fmt.Errorf("update article by id: article ID:%d ,title:%s ,content:%s\n%w: %v",ID,title,content ,repository.ErrDatabase,err)
	}

	return article, nil
}

func (a *articlePostgresRepo) DeleteByID(ctx context.Context, ID uint) error {
	article := new(model.Article)

	result := a.postgresCLI.WithContext(ctx).Delete(article, ID)
	if result.Error != nil {
		return fmt.Errorf("delete article by id: article ID:%d \n%w: %v",ID,repository.ErrDatabase,result.Error)
	}

	if result.RowsAffected == 0 {
		return repository.ErrArticleNotFound
	}

	return nil
}
