package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	postgres_repository "blog/database/postgres/repo"

	"blog/internal/service/authentication"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthHelper  auth_helper.AuthHeaderHelper
	AuthService authentication.AuthService
}

type (
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
)

func (h *AuthHandler) Login(ctx *gin.Context) {
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
		user, accessToken, refreshToken, err := h.AuthService.Login(li.email, li.password)
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
func (h *AuthHandler) Register(ctx *gin.Context) {
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
		user, verifyEmailToken, err := h.AuthService.Register(
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
func (h *AuthHandler) Logout(ctx *gin.Context) {
	go func() {
		token, err := h.AuthHelper.GetHeader(ctx)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		err = h.AuthService.Logout(token)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		h.AuthHelper.DeleteHeader(ctx)
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user logouted!", map[string]interface{}{},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
