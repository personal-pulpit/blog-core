package article

import (
	"blog/internal/model"
	"blog/internal/repository"
)

type ArticleService interface {
	Create(title, content string, userId model.ID) (*model.Article, error)
	Update(id model.ID, tilte, contet string) (*model.Article, error)
	Delete(articleId model.ID) error
	GetAll() ([]*model.Article, error)
	GetArticleByTitle(title string) (*model.Article, error)
	GetArticleById(id model.ID) (*model.Article, error)
}

type articleService struct {
	articlePostgresRepo repository.ArticlePostgresRepository
}

func NewAricleService(repo repository.ArticlePostgresRepository)ArticleService{
	return &articleService{repo}
}

func (s *articleService) Create(title, content string, autherId model.ID) (*model.Article, error)  {
	articleModel := model.NewArticle(title,content,autherId)

	article ,err := s.articlePostgresRepo.Create(articleModel)

	if err != nil {
		return nil,err
	}

	return article,nil
}

func (s *articleService) Update(id model.ID, tilte, contet string) (*model.Article, error) {
	article ,err := s.articlePostgresRepo.UpdateByID(id,tilte,contet)

	if err != nil {
		return nil,err
	}

	return article,nil
}

func (s *articleService) Delete(articleId model.ID) error {
	err := s.articlePostgresRepo.DeleteByID(articleId)

	if err != nil {
		return err
	}

	return nil
}

func (s *articleService)GetAll() ([]*model.Article, error){
	articles,err := s.articlePostgresRepo.GetAll()

	if err != nil{
		return nil,err
	}

	return articles,err
}

func (s *articleService)GetArticleByTitle(title string) (*model.Article, error){
	article,err := s.articlePostgresRepo.GetArticleByTitle(title)

	if err != nil{
		return nil,err
	}

	return article,err
}

func (s *articleService)GetArticleById(id model.ID) (*model.Article, error){
	article,err := s.articlePostgresRepo.GetArticleById(id)

	if err != nil{
		return nil,err
	}

	return article,err
}