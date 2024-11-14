package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ArticleTestSuite struct {
	suite.Suite
	repo    repository.ArticlePostgresRepository
	article *model.Article
}

func (s *ArticleTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewArticlePostgresRepo(db)
}

func (s *ArticleTestSuite) TestA_CreateArticle() {
	testCases := []struct {
		article *model.Article
		Valid   bool
	}{
		{
			article: model.NewArticle("article1", "content", "1"),
			Valid:   true,
		},
	}

	for _, tc := range testCases {
		article, err := s.repo.Create(tc.article)
		if tc.Valid {
			s.NoError(err)
			s.article = article
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(article)
		}
	}
}

func (s *ArticleTestSuite) TestB_GetAllArticles() {
	articles, err := s.repo.GetAll()
	s.NoError(err)
	s.NotNil(articles)
}

func (s *ArticleTestSuite) TestC_GetArticlesByTitle() {
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
		//fix it
		articles, err := s.repo.GetArticleByTitle(tc.Title)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(articles)
			// for _, article := range articles {
			// 	s.NotNil(article.Genres)
			// }
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(articles)
		}
	}
}

func (s *ArticleTestSuite) TestD_GetArticleByID() {
	testCases := []struct {
		ID    model.ID
		Valid bool
	}{
		{
			ID:    s.article.ID,
			Valid: true,
		},
		{
			ID:    "invalid",
			Valid: false,
		},
	}

	for _, tc := range testCases {
		article, err := s.repo.GetArticleById(tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(article)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(article)
		}
	}
}

func (s *ArticleTestSuite) TestE_ModifyArticle() {
	testCases := []struct {
		articleID model.ID
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
			articleID: "invalid",
			title:     "article 2",
			content:   "my fav article",
			authorID:  "2",
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		savedArticle, err := s.repo.UpdateByID(tc.articleID, tc.title, tc.content)
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
	testCases := []struct {
		articleID model.ID
		Valid     bool
	}{
		{
			articleID: s.article.ID,
			Valid:     true,
		},

		{
			articleID: "invalid",
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		err := s.repo.DeleteByID(tc.articleID)
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
