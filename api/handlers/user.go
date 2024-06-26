package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	postgres_repository "blog/database/postgres/repo"

	"blog/internal/service/authentication"
	"blog/internal/service/user"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		AuthHelper  auth_helper.AuthHeaderHelper
		AuthService authentication.AuthService
		UserService user.UserService
	}
	updateInput struct {
		firstName string `form:"firstName" binding:"required"`
		lastName  string `form:"lastName" binding:"required"`
		biography string `form:"biography" binding:"required"`
	}
	signinInput struct {
		firstName string `form:"firstName" binding:"required"`
		lastName  string `form:"lastName" binding:"required"`
		password  string `form:"password" binding:"required"`
		email     string `form:"email" binding:"required,emailvalidatior"`
		biography string `form:"biography" binding:"required"`
	}
	loginInput struct {
		email    string `form:"email" binding:"required,emailvalidatior"`
		password string `form:"password" binding:"required"`
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

func (u *User) GetProfile(ctx *gin.Context) {
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
func (u *User) Login(ctx *gin.Context) {
	go func() {
		var li loginInput
		err := ctx.ShouldBind(&li)
		if err != nil {
			if utils.CheckErrorForWord(err, "required") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrPleaseCompleteAllFields),
					nil)
				return
			} else if utils.CheckErrorForWord(err, "usernamevalidaitor") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrUsernameShouldContain),
					nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		user, accessToken, refreshToken, err := u.AuthService.Login(li.email, li.password)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrUserNotFound) || errors.Is(err, postgres_repository.ErrUsernameOrPasswordWrong) {
				userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "welcome back!", map[string]interface{}{
				"user":          user,
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) Register(ctx *gin.Context) {
	go func() {
		var si signinInput
		err := ctx.ShouldBind(&si)
		if err != nil {
			if utils.CheckErrorForWord(err, "required") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrPleaseCompleteAllFields),
					nil)
				return
			} else if utils.CheckErrorForWord(err, "emailvalidatior") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrInvalidEmail),
					nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		user, verifyEmailToken, err := u.AuthService.Register(
			si.firstName,
			si.lastName,
			si.email,
			si.biography,
			si.password,
		)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrEmailAlreadyExits) ||
				errors.Is(err, postgres_repository.ErrUsernameAlreadyExits) ||
				errors.Is(err, postgres_repository.ErrPhoneNumberAlreadyExits) {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusCreated, "user created!", map[string]interface{}{
				"user":             user,
				"verifyEmailToken": verifyEmailToken,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusCreated, userResponseChannel)
}
func (u *User) UpdateProfile(ctx *gin.Context) {
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
func (u *User) DeleteAccount(ctx *gin.Context) {
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
func (u *User) Logout(ctx *gin.Context) {
	go func() {
		token, err := u.AuthHelper.GetHeader(ctx)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		err = u.AuthService.Logout(token)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		u.AuthHelper.DeleteHeader(ctx)
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user logouted!", map[string]interface{}{},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
