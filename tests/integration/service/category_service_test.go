package service

import (
	"blog/internal/service/category"

	postgres_repository "blog/database/postgres/repo"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CategoryTestSuite struct {
	suite.Suite

	service category.CategoryService

	categoryID uint
}

func (s *CategoryTestSuite) SetupSuite() {
	categoryRepo := postgres_repository.NewCategoryRepository(db)

	s.service = category.NewCategoryService(categoryRepo)
}

func (s *CategoryTestSuite) TestA_CreateCategory() {
	ctx := context.TODO()

	category, err := s.service.CreateCategory(ctx, "category1")

	s.NoError(err)
	s.NotNil(category)

	s.categoryID = category.ID
}

func (s *CategoryTestSuite) TestB_GetAllCategories() {
	ctx := context.TODO()

	categories, err := s.service.GetAllCategories(ctx)
	s.NoError(err)
	s.NotNil(categories)
}

func (s *CategoryTestSuite) TestC_GetCategoryByID() {
	ctx := context.TODO()

	testCases := []struct {
		categoryID uint
		Valid      bool
	}{
		{
			categoryID: 0,
			Valid:      false,
		},
		{
			categoryID: s.categoryID,
			Valid:      true,
		},
	}

	for _, tc := range testCases {
		category, err := s.service.GetCategoryByID(ctx, tc.categoryID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(category)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(category)
		}
	}
}

func (s *CategoryTestSuite) TestD_DeleteCategory() {
	ctx := context.TODO()

	testCases := []struct {
		categoryID uint
		Valid      bool
	}{
		{
			categoryID: s.categoryID,
			Valid:      true,
		},
		{
			categoryID: 0,
			Valid:      false,
		},
		{
			categoryID: s.categoryID,
			Valid:      false,
		},
	}

	for _, tc := range testCases {
		err := s.service.DeleteCategoryByID(ctx, s.categoryID)
		if tc.Valid {
			s.NoError(err)

		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestCategorySuite(t *testing.T) {
	t.Helper()

	suite.Run(t, new(CategoryTestSuite))
}
