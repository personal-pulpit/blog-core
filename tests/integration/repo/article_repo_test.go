package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ArticleTestSuite struct {
	suite.Suite
	repo        repository.ArticleRepository
	userRepo    repository.UserRepository
	categoryRepo repository.CategoryRepository
	commentRepo repository.CommentRepository
	article     *model.Article

	articleAuthorID uint
	articleCategory *model.Category
}

func (s *ArticleTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewArticlePostgresRepo(db)
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)
	s.commentRepo = postgres_repository.NewCommentPostgresRepository(db)
	categoryRepo := postgres_repository.NewCategoryRepository(db)

	user, tx, err := s.userRepo.Create(context.TODO(), model.NewUser("user1", "user1", "<EMAIL>", "user1"))
	s.Nil(err)
	s.NotNil(user)
	tx.Commit()

	s.articleAuthorID = user.ID

	category, err := categoryRepo.CreateCategory(context.TODO(), model.NewCategory("category1"))
	s.Nil(err)
	s.NotNil(category)

	s.categoryRepo = categoryRepo
	s.articleCategory = category
}

func (s *ArticleTestSuite) TestA_CreateArticle() {
	ctx := context.TODO()

	notExistingCategory1 := model.NewCategory("notExistingCategory")
	notExistingCategory2 := model.NewCategory("notExistingCategory2")
	notExistingCategory2.ID = 5

	testCases := []struct {
		article *model.Article
		Valid   bool
	}{
		{
			article: model.NewArticle("article1", "content", s.articleAuthorID, []model.Category{*s.articleCategory}),
			Valid:   true,
		},
		{
			article: model.NewArticle("article0", "content", s.articleAuthorID, []model.Category{*notExistingCategory1}),
			Valid:  false,
		},
		{
			article: model.NewArticle("article1", "content", s.articleAuthorID, []model.Category{*notExistingCategory2}),
			Valid:  false,
		},
		{
			article: model.NewArticle("article2", "content", 0, []model.Category{*s.articleCategory}),
			Valid:   false,
		},
	}

	for _, tc := range testCases {
		article, err := s.repo.Create(ctx, tc.article)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(article)
			s.article = article

			comment, err := s.commentRepo.Create(ctx, model.NewComment(article.ID, s.articleAuthorID, "comment1"))
			s.Nil(err)
			s.NotNil(comment)

			comment2, err := s.commentRepo.Create(ctx, model.NewComment(article.ID, s.articleAuthorID, "comment2"))
			s.Nil(err)
			s.NotNil(comment2)

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

func (s *ArticleTestSuite) TestC_SearchArticle() {
	ctx := context.TODO()
	fakeTitleData := "fake"
	fakeTimeNowData := time.Now().Format("2006-01-02 15:04:05.999999999-07")

	timeNow, err := utils.ParasTime(fakeTimeNowData)
	s.Nil(err)

	testCases := []struct {
		filterData *model.ArticleFilter
		Valid      bool
	}{
		{
			filterData: &model.ArticleFilter{Title: &s.article.Title, PublishedAt: timeNow},
			Valid:      true,
		},
		{
			filterData: &model.ArticleFilter{Title: &fakeTitleData, PublishedAtLT: timeNow},
			Valid:      true,
		},
		{
			filterData: &model.ArticleFilter{Title: &fakeTitleData, PublishedAt: &s.article.CreatedAt, PublishedAtGT: timeNow},
			Valid:      true,
		},
	}

	for _, tc := range testCases {
		articles, err := s.repo.SearchArticles(ctx, tc.filterData)
		if tc.Valid {
			s.NoError(err)
			for _, article := range articles {
				s.NotNil(article)
				s.NotNil(article.Author)
				s.NotNil(article.Categories)
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
			s.NotNil(article.Comments)
			s.NotNil(article.Categories)
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
			articleID: 0,
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
			s.Nil(savedArticle)
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

	// Ensure the article's comments is exists
	comments, err := s.commentRepo.GetAll(ctx)
	s.NoError(err)
	s.NotEqual(0, len(comments))

	for _, tc := range testCases {
		err := s.repo.DeleteByID(ctx, tc.articleID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}

	// Ensure the article is deleted
	_, err = s.repo.GetArticleByID(ctx, s.article.ID)
	s.Error(err)

	// Ensure the article's user exists
	user, err := s.userRepo.GetUserByID(ctx, s.articleAuthorID)
	s.Nil(err)
	s.NotNil(user)

	// Ensure the article's comments is deleted
	comments, err = s.commentRepo.GetAll(ctx)
	s.NoError(err)
	s.Equal(0, len(comments))

	// Ensure the article's category exists
	category, err := s.categoryRepo.GetCategoryByID(ctx,s.articleCategory.ID)
	s.NoError(err)
	s.NotNil(category)
}

func TestArticleTestSuite(t *testing.T) {
	suite.Run(t, new(ArticleTestSuite))
}
