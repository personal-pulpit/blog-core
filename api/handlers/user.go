package handlers

import (
	"blog/api/helpers"
	objStorage "blog/pkg/object_storage"
	"blog/utils/file"
	"time"

	"blog/internal/repository"
	"blog/internal/service/user"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService user.UserService
}

type (
	updateInput struct {
		FirstName string `json:"firstName" binding:"required"`
		LastName  string `json:"lastName" binding:"required"`
		Biography string `json:"biography" binding:"required"`
	}
	deleteAccountInput struct {
		Password string `json:"password" binding:"required"`
	}
)

var (
	ErrPleaseCompleteAllFields = errors.New("please complete all fields")
	ErrInvalidEmail            = errors.New("email is invalid")
)

// @Summary Get current user profile
// @Description Retrieve the profile of the currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile retrieved successfully"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "User not found"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to retrieve user profile"
// @Router /api/v1/users [get]
func (u *UserHandler) GetCurrentUser(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))

	user, err := u.UserService.GetUserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "User not found", map[string]interface{}{
				"error":   err.Error(),
				"user_id": id,
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve user profile", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
		"user": user,
		"role": ctx.GetString("role"),
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
		},
	})
}

// @Summary Get user profile by ID
// @Description Retrieve the profile of a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile retrieved successfully"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "User not found"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to retrieve user profile"
// @Router /api/v1/users/{id} [get]
func (u *UserHandler) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	ID, err := helpers.StringToInt(id)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid user ID", map[string]interface{}{
			"error":         "User ID must be a valid integer",
			"provided_data": id,
		})
		return
	}

	user, err := u.UserService.GetUserProfile(ctx, uint(ID))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			helpers.RespondWithError(ctx, http.StatusNotFound, "User not found", map[string]interface{}{
				"error":   err.Error(),
				"user_id": id,
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve user profile", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Profile retrieved successfully", map[string]interface{}{
		"user": user,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
		},
	})
}

// @Summary Get profile image URL
// @Description Retrieve the profile image URL of a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile image URL retrieved successfully"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile image not found"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to retrieve profile image URL"
// @Router /api/v1/users/{id}/profile-picture-url/ [get]
func (u *UserHandler) GetProfileImageURL(ctx *gin.Context) {
	id := ctx.Param("id")

	ID, err := helpers.StringToInt(id)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid user ID", map[string]interface{}{
			"error":         "User ID must be a valid integer",
			"provided_data": id,
		})
		return
	}

	url, err := u.UserService.GetProfileImageURL(ctx, uint(ID))
	if err != nil {
		if err == objStorage.ErrObjectNotFound {
			helpers.RespondWithError(ctx, http.StatusNotFound, "Failed to retrieve profile image URL with provided id", map[string]interface{}{
				"user_id": id,
				"error":   err.Error(),
			})

			return
		}

		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to retrieve profile image URL", map[string]interface{}{
			"error": err.Error(),
		})

		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Profile image URL retrieved successfully", map[string]interface{}{
		"url": url,
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
		},
	})
}

// @Summary Upload profile image
// @Description Upload a profile image for the currently authenticated user
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Profile image file"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile image uploaded successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid image file"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to upload profile image"
// @Router /api/v1/users/profile-picture [post]
func (u *UserHandler) AddProfileImage(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))

	imageFile, err := ctx.FormFile("image")
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to upload profile image", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fileData, err := file.CheckFile(imageFile)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Failed to upload profile image", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	err = u.UserService.AddProfileImage(ctx, id, fileData)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to upload profile image", map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Profile Image uploaded successfully", map[string]interface{}{
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
		},
	})
}

// @Summary Update user profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param user body updateInput true "User profile data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile updated successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid update data"
// @Router /api/v1/users [patch]
func (u *UserHandler) UpdateProfile(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))
	var ui updateInput
	err := ctx.ShouldBindJSON(&ui)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Missing required fields", map[string]interface{}{
				"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
				"required_fields":   []string{"firstName", "lastName", "biography"},
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid update data", map[string]interface{}{
			"validation_errors": utils.GetValidationError(err),
		})
		return
	}

	user, err := u.UserService.UpdateProfile(
		ctx,
		id,
		ui.FirstName,
		ui.LastName,
		ui.Biography,
	)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Profile update failed", map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Profile updated successfully", map[string]interface{}{
		"user": user,
		"metadata": map[string]interface{}{
			"timestamp":  time.Now(),
			"updated_at": time.Now(),
		},
	})
}

// @Summary Delete user account
// @Description Delete the account of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param input body deleteAccountInput true "Account deletion data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Account deleted successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid deletion request"
// @Router /api/v1/users [delete]
func (u *UserHandler) DeleteAccount(ctx *gin.Context) {
	var input deleteAccountInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			helpers.RespondWithError(ctx, http.StatusBadRequest, "Password is required for account deletion", map[string]interface{}{
				"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
				"required_fields":   []string{"password"},
			})
			return
		}
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Invalid deletion request", map[string]interface{}{
			"validation_errors": utils.GetValidationError(err),
		})
		return
	}

	id := uint(ctx.GetInt("id"))
	err = u.UserService.DeleteAccount(ctx, id, input.Password)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusBadRequest, "Account deletion failed", map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Account deleted successfully", map[string]interface{}{
		"metadata": map[string]interface{}{
			"deleted_at": time.Now(),
			"user_id":    id,
		},
	})
}

// @Summary Delete profile image
// @Description Delete the profile image of the currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile image deleted successfully"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to delete profile image"
// @Router /api/v1/users/profile-picture [delete]
func (u *UserHandler) DeleteProfileImage(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))

	err := u.UserService.DeleteProfileImage(ctx, id)
	if err != nil {
		helpers.RespondWithError(ctx, http.StatusInternalServerError, "Failed to delete profile image", map[string]interface{}{
			"error":   err.Error(),
			"user_id": id,
		})
		return
	}

	helpers.RespondWithSuccess(ctx, http.StatusOK, "Profile Image deleted successfully", map[string]interface{}{
		"metadata": map[string]interface{}{
			"timestamp": time.Now(),
		},
	})
}
