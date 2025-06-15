package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BookmarkTestSuite struct {
	suite.Suite
	repo        repository.BookmarkRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	bookmarkID  uint
	articleID   uint
	userID      uint
}

func (s *BookmarkTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewBookmarkPostgresRepository(db)
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)
	s.articleRepo = postgres_repository.NewArticlePostgresRepo(db)
	categoryRepo := postgres_repository.NewCategoryRepository(db)

	user, tx, err := s.userRepo.Create(context.TODO(), model.NewUser("user2", "user2", "<TEST-BOOKMARK>", "user2"))
	s.Nil(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	category, err := categoryRepo.CreateCategory(context.TODO(), model.NewCategory("bookmark-category1"))
	s.Nil(err)
	s.NotNil(category)

	article, err := s.articleRepo.Create(context.TODO(), model.NewArticle("article1", "content", s.userID, []model.Category{*category}))
	s.Nil(err)
	s.NotNil(article)

	s.articleID = article.ID
}

func (s *BookmarkTestSuite) TestA_CreateBookmark() {
	ctx := context.TODO()

	testCases := []struct {
		bookmark *model.Bookmark
		Valid    bool
	}{
		{
			bookmark: model.NewBookmark(s.userID, s.articleID),
			Valid:    true,
		},
		{
			bookmark: model.NewBookmark(0, s.articleID),
			Valid:    false,
		},
		{
			bookmark: model.NewBookmark(s.userID, 0),
			Valid:    false,
		},
	}

	for _, tc := range testCases {
		bookmark, err := s.repo.Create(ctx, tc.bookmark)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(bookmark)
			s.bookmarkID = bookmark.ID

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(bookmark)
		}
	}
}

func (s *BookmarkTestSuite) TestB_GetUserBookmarks() {
	ctx := context.TODO()


	bookmarks, err := s.repo.GetByUserID(ctx, s.userID)
	s.NoError(err)
	s.NotNil(bookmarks)
	
}


func (s *BookmarkTestSuite) TestC_GetBookmarkByID() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    s.bookmarkID,
			Valid: true,
		},
		{
			ID:    0,
			Valid: false,
		},
	}

	for _, tc := range testCases {
		bookmark, err := s.repo.GetByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(bookmark)
			s.NotNil(bookmark.Article)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(bookmark)
		}
	}
}


func (s *BookmarkTestSuite) TestD_DeleteBookmarkByArticleIDAndUserID() {
	ctx := context.TODO()

	testCases := []struct {
		articleID uint
		userID    uint
		Valid     bool
	}{
		{
			articleID: s.articleID,
			userID:    s.userID,
			Valid:     true,
		},

		{
			articleID: s.articleID,
			userID:    0,
			Valid:     false,
		},

		{
			articleID: 0,
			userID:    s.userID,
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		err := s.repo.DeleteByArticleIDAndUserID(ctx, tc.articleID,tc.userID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}

	//Ensure user exists after deleting its bookmark
	user, err := s.userRepo.GetUserByID(ctx, s.userID)
	s.Nil(err)
	s.NotNil(user)

	//Ensure article exists after deleting its bookmark
	article, err := s.articleRepo.GetArticleByID(ctx, s.articleID)
	s.Nil(err)
	s.NotNil(article)
}

func TestBookmarkTestSuite(t *testing.T) {
	suite.Run(t, new(BookmarkTestSuite))
}
