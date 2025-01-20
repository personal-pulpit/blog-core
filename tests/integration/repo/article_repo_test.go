package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ArticleTestSuite struct {
	suite.Suite
	repo            repository.ArticlePostgresRepository
	userRepo        repository.UserPostgresRepository
	article         *model.Article
	articleAuthorID uint
}

func (s *ArticleTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewArticlePostgresRepo(db)
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)

	user, tx, err := s.userRepo.Create(context.TODO(), model.NewUser("user1", "user1", "<EMAIL>", "user1"))
	s.Nil(err)
	s.NotNil(user)
	tx.Commit()
	s.articleAuthorID = user.ID
}

func (s *ArticleTestSuite) TestA_CreateArticle() {
	ctx := context.TODO()

	testCases := []struct {
		article *model.Article
		Valid   bool
	}{
		{
			article: model.NewArticle("article1", "content", s.articleAuthorID),
			Valid:   true,
		},
		{
			article: model.NewArticle("article1", "content", 0),
			Valid:   false,
		},

	}

	for _, tc := range testCases {
		article, err := s.repo.Create(ctx, tc.article)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(article)
			s.article = article
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(article)
		}
	}
}

func (s *ArticleTestSuite) TestB_GetAllArticles() {
	ctx := context.TODO()

	articles, err := s.repo.GetAll(ctx)
	s.NoError(err)
	s.NotNil(articles)
	for _, article := range articles {
		s.NotNil(article)
		s.NotNil(article.Author)
	}
}

func (s *ArticleTestSuite) TestC_GetArticlesByTitle() {
	ctx := context.TODO()

	testCases := []struct {
		Title string
		Valid bool
	}{
		{
			Title: s.article.Title,
			Valid: true,
		},
		{
			Title: "Invalid",
			Valid: false,
		},
	}

	for _, tc := range testCases {
		articles, err := s.repo.GetArticleByTitle(ctx, tc.Title)
		if tc.Valid {
			s.NoError(err)
			for _, article := range articles {
				s.NotNil(article)
				s.NotNil(article.Author)
			}

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(articles)
		}
	}
}

func (s *ArticleTestSuite) TestD_GetArticleByID() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    s.article.ID,
			Valid: true,
		},
		{
			ID:    1000,
			Valid: false,
		},
	}

	for _, tc := range testCases {
		article, err := s.repo.GetArticleByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(article)
			s.NotNil(article.Author)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(article)
		}
	}
}

func (s *ArticleTestSuite) TestE_ModifyArticle() {
	ctx := context.TODO()

	testCases := []struct {
		articleID uint
		title     string
		content   string
		authorID  string
		Valid     bool
	}{
		{
			articleID: s.article.ID,
			title:     "article 2",
			content:   "my fav article",
			authorID:  "2",
			Valid:     true,
		},

		{
			articleID: 1000,
			title:     "article 2",
			content:   "my fav article",
			authorID:  "2",
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		savedArticle, err := s.repo.UpdateByID(ctx, tc.articleID, tc.title, tc.content)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(savedArticle)
			s.NotEmpty(savedArticle.UpdatedAt)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *ArticleTestSuite) TestF_DeleteArticleByID() {
	ctx := context.TODO()

	testCases := []struct {
		articleID uint
		Valid     bool
	}{
		{
			articleID: s.article.ID,
			Valid:     true,
		},

		{
			articleID: 1000,
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		err := s.repo.DeleteByID(ctx, tc.articleID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
