package service

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/internal/service/bookmark"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BookmarkTestSuite struct {
	suite.Suite
	bookmarkRepo repository.BookmarkRepository

	service bookmark.BookmarkService

	bookmarkID uint
	articleID uint
	userID    uint
}

func (s *BookmarkTestSuite) SetupSuite() {
	s.bookmarkRepo = postgres_repository.NewBookmarkPostgresRepository(db)
	userRepo := postgres_repository.NewUserPostgresRepository(db)
	articleRepo := postgres_repository.NewArticlePostgresRepo(db)
	categoryRepo := postgres_repository.NewCategoryRepository(db)

	userModel := model.NewUser("user1 firstName", "user1 lastName", "qazwsxedcrfv@fake.come", "user1 biography")
	user, tx, err := userRepo.Create(context.TODO(), userModel)
	s.NoError(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	category, err := categoryRepo.CreateCategory(context.TODO(), model.NewCategory("bookmark-category1-service"))
	s.NoError(err)
	s.NotNil(category)


	article, err := articleRepo.Create(context.TODO(), model.NewArticle("article1", "content of article1", user.ID, []model.Category{*category}))
	s.NoError(err)
	s.NotNil(article)

	s.articleID = article.ID

	s.service = bookmark.NewBookmarkService(s.bookmarkRepo, userRepo)
}

func (s *BookmarkTestSuite) TestA_CreateBookmark() {
	ctx := context.TODO()

	bookmark, err := s.service.CreateBookmark(ctx, s.userID, s.articleID)

	s.NoError(err)
	s.NotNil(bookmark)

	s.bookmarkID = bookmark.ID
}

func (s *BookmarkTestSuite) TestB_GetBookmarkByID() {
	ctx := context.TODO()

	testCases := []struct {
		bookmarkID uint
		Valid      bool
	}{
		{
			bookmarkID: 0,
			Valid:      false,
		},
		{
			bookmarkID: s.bookmarkID,
			Valid:      true,
		},
	}

	for _, tc := range testCases {
		bookmark, err := s.service.GetBookmarkByID(ctx, tc.bookmarkID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(bookmark)
		} else {
			s.Error(err)
			s.Nil(bookmark)
		}
	}
}

func (s *BookmarkTestSuite) TestC_GetUserBookmarks(){
	ctx := context.TODO()

	testCases := []struct {
		userID    uint
		Valid     bool
	}{
		{
			userID:    s.userID,
			Valid:     true,
		},
		{
			userID:    0,
			Valid:     false,
		},

	}

	for _, tc := range testCases {
		bookmarks,err := s.service.GetUsersBookmarks(ctx, tc.userID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(bookmarks)

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(bookmarks)
		}
	}
}



func (s *BookmarkTestSuite) TestD_DeleteBookmark() {
	ctx := context.TODO()

	testCases := []struct {
		userID    uint
		articleID uint
		Valid     bool
	}{
		{
			userID:    s.userID,
			articleID:0,
			Valid:     false,
		},
		{
			userID: 0,
			articleID: s.articleID,
			Valid:     false,
		},
		{
			userID:    s.userID,
			articleID: s.articleID,
			Valid:     true,
		},
	}

	
	for _, tc := range testCases {
		err := s.service.DeleteBookmark(ctx, tc.articleID, tc.userID)
		if tc.Valid {
			s.NoError(err)

		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestBookmarkSuite(t *testing.T) {
	t.Helper()

	suite.Run(t, new(BookmarkTestSuite))
}
