package article

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"context"
	"fmt"
)

type ArticleService interface {
	Create(ctx context.Context, title, content string, authorID uint, categoriesName []string) (*model.Article, error)
	Update(ctx context.Context, ID uint, title, content string) (*model.Article, error)
	Delete(ctx context.Context, articleID uint) error
	GetAll(ctx context.Context) ([]model.Article, error)
	SearchArticle(ctx context.Context, filter map[string]string) ([]model.Article, error)
	GetArticleByID(ctx context.Context, ID uint) (*model.Article, error)
}

type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
}

func NewArticleService(articleRepo repository.ArticleRepository, categoryRepo repository.CategoryRepository) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *articleService) Create(ctx context.Context, title, content string, authorID uint, categoriesName []string) (*model.Article, error) {
	var categories = []model.Category{}

	for _, categoryName := range categoriesName {
		category, err := s.categoryRepo.GetCategoryByName(ctx, categoryName)
		if err != nil {
			return nil, err
		}

		categories = append(categories, *category)
	}

	articleModel := model.NewArticle(title, content, authorID, categories)

	article, err := s.articleRepo.Create(ctx, articleModel)

	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleService) Update(ctx context.Context, ID uint, title, content string) (*model.Article, error) {
	article, err := s.articleRepo.UpdateByID(ctx, ID, title, content)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleService) Delete(ctx context.Context, articleID uint) error {
	err := s.articleRepo.DeleteByID(ctx, articleID)
	if err != nil {
		return err
	}

	return nil
}

func (s *articleService) GetAll(ctx context.Context) ([]model.Article, error) {
	articles, err := s.articleRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return articles, err
}

func (s *articleService) SearchArticle(ctx context.Context, filter map[string]string) ([]model.Article, error) {
	var articleFilter = &model.ArticleFilter{}

	for i, v := range filter {
		if i == "title" {
			articleFilter.Title = &v
		} else if i == "publishedAt" {
			parsedTime, err := utils.ParasTime(v)
			if err != nil {
				return nil, fmt.Errorf("SearchArticle:%w: %v : \ntime: %s", ErrParseTime, err, v)
			}

			articleFilter.PublishedAt = parsedTime
		} else if i == "publishedAtGt" {
			parsedTime, err := utils.ParasTime(v)
			if err != nil {
				return nil, fmt.Errorf("SearchArticle:%w: %v : \ntime: %s", ErrParseTime, err, v)
			}

			articleFilter.PublishedAtGT = parsedTime
		} else if i == "publishedAtLt" {
			parsedTime, err := utils.ParasTime(v)
			if err != nil {
				return nil, fmt.Errorf("SearchArticle:%w: %v : \ntime: %s", ErrParseTime, err, v)
			}

			articleFilter.PublishedAtLT = parsedTime
		}
	}

	articles, err := s.articleRepo.SearchArticles(ctx, articleFilter)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (s *articleService) GetArticleByID(ctx context.Context, ID uint) (*model.Article, error) {
	article, err := s.articleRepo.GetArticleByID(ctx, ID)

	if err != nil {
		return nil, err
	}

	return article, err
}
