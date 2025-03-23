package user

import (
	"blog/internal/model"
	"blog/internal/repository"
	objStorage "blog/pkg/object_storage"
	"blog/utils/file"
	"context"
)

type UserService interface {
	AddProfileImage(ctx context.Context, ID uint, image []byte) error
	GetUserProfile(ctx context.Context, ID uint) (*model.User, error)
	GetProfileImageURL(ctx context.Context, ID uint) (string, error)
	UpdateProfile(ctx context.Context, ID uint, FirstName, lastName, biography string) (*model.User, error)
	DeleteAccount(ctx context.Context, ID uint, password string) error
	DeleteProfileImage(ctx context.Context, ID uint) error
}

type userManager struct {
	userPostgresRepo  repository.UserRepository
	userCacheRepo     repository.UserCacheRepository
	authPostgresRepo  repository.AuthRepository
	objectStorageRepo objStorage.StorageRepository
}

func NewUserService(userPostgresRepo repository.UserRepository,userCacheRepo repository.UserCacheRepository, authPostgresRepo repository.AuthRepository, objectStorage objStorage.StorageRepository) UserService {
	return &userManager{
		userPostgresRepo:  userPostgresRepo,
		authPostgresRepo:  authPostgresRepo,
		objectStorageRepo: objectStorage,
		userCacheRepo:     userCacheRepo,
	}
}

func (u *userManager) AddProfileImage(ctx context.Context, ID uint, image []byte) error {
	fileName := file.GenerateFileName(ID)

	object := objStorage.Object{
		Name:        fileName,
		Data:        image,
		ContentType: objStorage.ImageContentType,
		Size:        int64(len(image)),
	}

	err := u.objectStorageRepo.UploadFile(ctx, objStorage.ProfileImageBucketName, object)
	if err != nil {
		return err
	}

	url,err := u.GetProfileImageURL(ctx, ID)
	if err != nil {
		return err
	}

	err = u.userCacheRepo.SetUserProfileImageURL(ctx, ID, url,objStorage.ExpirationTime)
	if err != nil {
		return err
	}

	return nil
}

func (u *userManager) GetUserProfile(ctx context.Context, ID uint) (*model.User, error) {
	userModel, err := u.userPostgresRepo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (u *userManager) GetProfileImageURL(ctx context.Context, ID uint) (string, error) {
	cachedURL,err := u.userCacheRepo.GetUserProfileImageURL(ctx,ID)
	
	//TODO log error
	if err == nil {
		return cachedURL, nil
	}

	fileName := file.GenerateFileName(ID)

	url, err := u.objectStorageRepo.GetObjectURL(ctx, objStorage.ProfileImageBucketName, fileName)
	if err != nil {
		return "", err
	}

	err = u.userCacheRepo.SetUserProfileImageURL(ctx, ID, url,objStorage.ExpirationTime )
	if err != nil {
		return "", err
	}

	return url, nil
}

func (u *userManager) UpdateProfile(ctx context.Context, ID uint, FirstName, lastName, biography string) (*model.User, error) {
	user, err := u.userPostgresRepo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	user.FirstName = FirstName
	user.LastName = lastName
	user.Biography = biography

	userModel, err := u.userPostgresRepo.UpdateByID(ctx, ID, user.FirstName, user.LastName, user.Biography)
	if err != nil {
		return nil, err
	}
	return userModel, nil
}

func (u *userManager) DeleteProfileImage(ctx context.Context, ID uint) error {
	fileName := file.GenerateFileName(ID)

	err := u.objectStorageRepo.DeleteFile(ctx, objStorage.ProfileImageBucketName, fileName)
	if err != nil {
		return err
	}

	err = u.userCacheRepo.DestroyUserProfileImageURL(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userManager) DeleteAccount(ctx context.Context, ID uint, password string) error {
	auth, err := u.authPostgresRepo.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	if !auth.EmailVerified {
		return ErrUsersEmailNotVerified
	}

	err = u.authPostgresRepo.DeleteByID(ctx, ID)
	if err != nil {
		return err
	}

	err = u.userPostgresRepo.DeleteByID(ctx, ID)
	if err != nil {
		return err
	}

	return nil
}
