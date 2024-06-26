package handlers

import (
	"blog/api/helpers"
	postgres_repository "blog/database/postgres/repo"

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
		firstName string `form:"firstName" binding:"required"`
		lastName  string `form:"lastName" binding:"required"`
		biography string `form:"biography" binding:"required"`
	}
	deleteAccountInput struct {
		password string `form:"password" binding:"required"`
	}
)

var (
	userResponseChannel        = make(chan helpers.HttpResponse)
	ErrPleaseCompleteAllFields = errors.New("please complete all fields")
	ErrUsernameShouldContain   = errors.New("username should contain: a-z  _ 0-9")
	ErrInvalidEmail            = errors.New("email is invalid")
	ErrInvalidPhonenumber      = errors.New("phonenumber is invalid")
)

func (u *UserHandler) GetProfile(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		user, err := u.UserService.GetUserProfile(id)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrUserNotFound) {
				userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusCreated, "user Got!", map[string]interface{}{
				"user": user,
			},
		)

	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *UserHandler) UpdateProfile(ctx *gin.Context) {
	go func() {
		id := ctx.GetString("id")
		var ui updateInput
		err := ctx.ShouldBind(&ui)
		if err != nil {
			if utils.CheckErrorForWord(err, "required") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrPleaseCompleteAllFields),
					nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		user, err := u.UserService.UpdateProfile(
			id,
			ui.firstName,
			ui.lastName,
			ui.biography,
		)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user updated!", map[string]interface{}{
				"user": user,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *UserHandler) DeleteAccount(ctx *gin.Context) {
	go func() {
		id := ctx.GetString("id")
		var si deleteAccountInput
		err := ctx.ShouldBind(&si)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		err = u.UserService.DeleteAccount(id, si.password)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrUserNotFound) {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user deleted!", map[string]interface{}{},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
