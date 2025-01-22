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
		FirstName string `form:"firstName" binding:"required"`
		LastName  string `form:"lastName" binding:"required"`
		Biography string `form:"biography" binding:"required"`
	}
	deleteAccountInput struct {
		Password string `form:"password" binding:"required"`
	}
)

var (
	ErrPleaseCompleteAllFields = errors.New("please complete all fields")
	ErrUsernameShouldContain   = errors.New("username should contain: a-z  _ 0-9")
	ErrInvalidEmail            = errors.New("email is invalid")
	ErrInvalidPhoneNumber      = errors.New("phone number is invalid")
)

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
			"metadata": map[string]interface{}{
				"timestamp": time.Now(),
			},
		}))
}

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

func (u *UserHandler) UpdateProfile(ctx *gin.Context) {
	id := uint(ctx.GetInt("id"))
	var ui updateInput
	err := ctx.ShouldBind(&ui)
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

func (u *UserHandler) DeleteAccount(ctx *gin.Context) {
	var input deleteAccountInput
	err := ctx.ShouldBind(&input)
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
