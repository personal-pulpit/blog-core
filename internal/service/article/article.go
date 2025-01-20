package article

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
)

type ArticleService interface {
	Create(ctx context.Context,title, content string, authorID uint) (*model.Article, error)
	Update(ctx context.Context,ID uint, title, content string) (*model.Article, error)
	Delete(ctx context.Context,articleID uint) error
	GetAll(ctx context.Context,) ([]model.Article, error)
	GetArticleByTitle(ctx context.Context,title string) ([]model.Article, error)
	GetArticleByID(ctx context.Context,ID uint) (*model.Article, error)
}

type articleService struct {
	articlePostgresRepo repository.ArticlePostgresRepository
}

func NewArticleService(repo repository.ArticlePostgresRepository)ArticleService{
	return &articleService{repo}
}

func (s *articleService) Create(ctx context.Context,title, content string, authorID uint) (*model.Article, error)  {
	articleModel := model.NewArticle(title,content,authorID)

	article ,err := s.articlePostgresRepo.Create(ctx,articleModel)

	if err != nil {
		return nil,err
	}

	return article,nil
}

func (s *articleService) Update(ctx context.Context,ID uint, title, content string) (*model.Article, error) {
	article ,err := s.articlePostgresRepo.UpdateByID(ctx,ID,title,content)
	if err != nil {
		return nil,err
	}

	return article,nil
}

func (s *articleService) Delete(ctx context.Context,articleID uint) error {
	err := s.articlePostgresRepo.DeleteByID(ctx,articleID)
	if err != nil {
		return err
	}

	return nil
}

func (s *articleService)GetAll(ctx context.Context) ([]model.Article, error){
	articles,err := s.articlePostgresRepo.GetAll(ctx)
	if err != nil{
		return nil,err
	}

	return articles,err
}

func (s *articleService)GetArticleByTitle(ctx context.Context,title string) ([]model.Article, error){
	articles,err := s.articlePostgresRepo.GetArticleByTitle(ctx,title)
	if err != nil{
		return nil,err
	}

	return articles,err
}

func (s *articleService)GetArticleByID(ctx context.Context,ID uint) (*model.Article, error){
	article,err := s.articlePostgresRepo.GetArticleByID(ctx,ID)

	if err != nil{
		return nil,err
	}

	return article,err
}