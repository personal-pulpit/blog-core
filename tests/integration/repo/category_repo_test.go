package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CategoryTestSuite struct {
	suite.Suite
	repo        repository.CategoryRepository
	articleRepo repository.ArticleRepository
	categoryID  uint
	articleID   uint
	userID      uint
}

func (s *CategoryTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewCategoryRepository(db)
	s.articleRepo = postgres_repository.NewArticlePostgresRepo(db)
	userRepo := postgres_repository.NewUserPostgresRepository(db)

	user, tx, err := userRepo.Create(context.TODO(), model.NewUser("user2", "user2", "<TEST>", "user2"))
	s.Nil(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID
}

func (s *CategoryTestSuite) TestA_CreateCategory() {
	ctx := context.TODO()

	testCases := []struct {
		Category *model.Category
		Valid    bool
	}{
		{
			Category: model.NewCategory("Category1"),
			Valid:    true,
		},
		{
			Category: model.NewCategory("Category1"),
			Valid:    false,
		},
	}

	for _, tc := range testCases {
		Category, err := s.repo.CreateCategory(ctx, tc.Category)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(Category)

			article, err := s.articleRepo.Create(ctx, model.NewArticle("articleCategoryTest", "This is category test", s.userID, []model.Category{*Category}))
			s.NoError(err)
			s.NotNil(article)

			s.articleID = article.ID
			s.categoryID = Category.ID

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(Category)
		}
	}
}

func (s *CategoryTestSuite) TestB_GetAllCategories() {
	ctx := context.TODO()

	Categories, err := s.repo.GetAllCategories(ctx)
	s.NoError(err)
	s.NotNil(Categories)
}

func (s *CategoryTestSuite) TestC_GetCategoryByID() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    s.categoryID,
			Valid: true,
		},
		{
			ID:    0,
			Valid: false,
		},
	}

	for _, tc := range testCases {
		category, err := s.repo.GetCategoryByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(category)
			s.NotNil(category.Articles)

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(category)
		}
	}
}

func (s *CategoryTestSuite) TestD_DeleteCategoryByID() {
	ctx := context.TODO()

	testCases := []struct {
		CategoryID uint
		Valid      bool
	}{
		{
			CategoryID: s.categoryID,
			Valid:      true,
		},

		{
			CategoryID: 0,
			Valid:      false,
		},

		{
			CategoryID: s.categoryID,
			Valid:      false,
		},
	}

	for _, tc := range testCases {
		err := s.repo.DeleteCategoryByID(ctx, tc.CategoryID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}

	//Ensure article exists after deleting its Category
	article, err := s.articleRepo.GetArticleByID(ctx, s.articleID)
	s.Nil(err)
	s.NotNil(article)
}

func TestCategoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
}
