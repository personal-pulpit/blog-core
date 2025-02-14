package handlers

import (
	"blog/api/helpers"
	postgres_repository "blog/database/postgres/repo"
	"time"

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
// @Router /api/v1/user/me [get]
func (u *UserHandler) GetCurrentUser(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))

	user, err := u.UserService.GetUserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, postgres_repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"User not found",
				map[string]interface{}{
					"error":   err.Error(),
					"user_id": id,
				}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Failed to retrieve user profile",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Profile retrieved successfully",
		map[string]interface{}{
			"user": user,
			"role": ctx.GetString("role"),
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

// @Summary Get user profile by ID
// @Description Retrieve the profile of a user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile retrieved successfully"
// @Failure 404 {object} helpers.HttpResponse{data=map[string]interface{}} "User not found"
// @Failure 500 {object} helpers.HttpResponse{data=map[string]interface{}} "Failed to retrieve user profile"
// @Router /api/v1/user/{id} [get]
func (u *UserHandler) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := u.UserService.GetUserProfile(ctx, uint(helpers.StringToInt(id)))
	if err != nil {
		if errors.Is(err, postgres_repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, helpers.NewHttpResponse(
				http.StatusNotFound,
				"User not found",
				map[string]interface{}{
					"error":   err.Error(),
					"user_id": id,
				}))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewHttpResponse(
			http.StatusInternalServerError,
			"Failed to retrieve user profile",
			map[string]interface{}{
				"error": err.Error(),
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Profile retrieved successfully",
		map[string]interface{}{
			"user": user,
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

// @Summary Update user profile
// @Description Update the profile of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param user body updateInput true "User profile data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Profile updated successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid update data"
// @Router /api/v1/user/update [patch]
func (u *UserHandler) UpdateProfile(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))
	var ui updateInput
	err := ctx.ShouldBindJSON(&ui)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Missing required fields",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"required_fields":   []string{"firstName", "lastName", "biography"},
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid update data",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
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
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Profile update failed",
			map[string]interface{}{
				"error":   err.Error(),
				"user_id": id,
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Profile updated successfully",
		map[string]interface{}{
			"user": user,
			"metadata": map[string]interface{}{
				"timestamp":  time.Now(),
				"updated_at": time.Now(),
			},
		}))
}

// @Summary Delete user account
// @Description Delete the account of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Param input body deleteAccountInput true "Account deletion data"
// @Success 200 {object} helpers.HttpResponse{data=map[string]interface{}} "Account deleted successfully"
// @Failure 400 {object} helpers.HttpResponse{data=map[string]interface{}} "Invalid deletion request"
// @Router /api/v1/user/delete [delete]
func (u *UserHandler) DeleteAccount(ctx *gin.Context) {
	var input deleteAccountInput
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		if utils.CheckErrorForWord(err, "required") {
			ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
				http.StatusBadRequest,
				"Password is required for account deletion",
				map[string]interface{}{
					"validation_errors": utils.GetValidationError(ErrPleaseCompleteAllFields),
					"required_fields":   []string{"password"},
				}))
			return
		}
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Invalid deletion request",
			map[string]interface{}{
				"validation_errors": utils.GetValidationError(err),
			}))
		return
	}

	id := uint(ctx.GetInt("id"))
	err = u.UserService.DeleteAccount(ctx, id, input.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewHttpResponse(
			http.StatusBadRequest,
			"Account deletion failed",
			map[string]interface{}{
				"error":   err.Error(),
				"user_id": id,
			}))
		return
	}

	ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
		http.StatusOK,
		"Account deleted successfully",
		map[string]interface{}{
			"metadata": map[string]interface{}{
				"deleted_at": time.Now(),
				"user_id":    id,
			},
		}))
}
