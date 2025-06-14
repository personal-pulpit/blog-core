package bookmark

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
)

type BookmarkService interface {
	CreateBookmark(ctx context.Context,userID, articleID uint) (*model.Bookmark, error)
	GetBookmarkByID(ctx context.Context,bookmarkID uint) (*model.Bookmark, error)
	GetUsersBookmarks(ctx context.Context,userID uint) ([]model.Bookmark, error)
	DeleteBookmark(ctx context.Context,ownerID,bookmarkID uint) error
}

type bookmarkServiceImpl struct {
	bookmarkRepo repository.BookmarkRepository
	userRepo    repository.UserRepository
}

func NewBookmarkService(bookmarkRepo repository.BookmarkRepository, userRepo repository.UserRepository) BookmarkService {
	return &bookmarkServiceImpl{
		bookmarkRepo: bookmarkRepo,
		userRepo:    userRepo,
	}
}

func (s *bookmarkServiceImpl) CreateBookmark(ctx context.Context,userID, articleID uint) (*model.Bookmark, error) {
	bookmark := model.NewBookmark(userID,articleID)

	return s.bookmarkRepo.Create(ctx, bookmark)
}

func (s *bookmarkServiceImpl) GetBookmarkByID(ctx context.Context, bookmarkID uint) (*model.Bookmark, error) {
	return s.bookmarkRepo.GetByID(ctx, bookmarkID)
}

func (s *bookmarkServiceImpl) GetUsersBookmarks(ctx context.Context, userID uint) ([]model.Bookmark, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.Bookmarks, nil
}


func (s *bookmarkServiceImpl) DeleteBookmark(ctx context.Context,ownerID, bookmarkID uint) error {
	err := s.checkBookmarkOwner(ctx, ownerID, bookmarkID)
	if err != nil {
		return err
	}

	return s.bookmarkRepo.DeleteByID(ctx, bookmarkID)
}


func (s *bookmarkServiceImpl) checkBookmarkOwner(ctx context.Context, ownerID, bookmarkID uint) error {
	bookmark, err := s.bookmarkRepo.GetByID(ctx, bookmarkID)
	if err != nil {
		return err
	}

	if bookmark.UserID != ownerID {
		return ErrPermissionDenied
	}
	return nil
}
