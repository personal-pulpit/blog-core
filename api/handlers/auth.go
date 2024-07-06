package handlers

import (
	"blog/api/helpers"
	postgres_repository "blog/database/postgres/repo"

	"blog/internal/service/authentication"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService authentication.AuthService
}

type (
	verifyEmailInput struct {
		OTP string `form:"otp" binding:"required"`
	}
	signinInput struct {
		FirstName string `form:"firstName" binding:"required"`
		LastName  string `form:"lastName" binding:"required"`
		Password  string `form:"password" binding:"required"`
		Email     string `form:"email" binding:"required,emailvalidatior"`
		Biography string `form:"biography" binding:"required"`
	}
	loginInput struct {
		Email    string `form:"email" binding:"required,emailvalidatior"`
		Password string `form:"password" binding:"required"`
	}
)

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
			si.FirstName,
			si.LastName,
			si.Email,
			si.Biography,
			si.Password,
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

func (h *AuthHandler) VerifyEmail(ctx *gin.Context) {
	go func() {
		var input verifyEmailInput

		err := ctx.ShouldBind(&input)
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
		err = h.AuthService.VerifyEmail(input.OTP, ctx.GetString("id"))
		if err != nil {
			if errors.Is(err, postgres_repository.ErrUserNotFound) || errors.Is(err, postgres_repository.ErrEmailOrPasswordWrong) {
				userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "welcome back!", map[string]interface{}{
				"emailVerified": true,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}

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
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		user, accessToken, refreshToken, err := h.AuthService.Login(li.Email, li.Password)
		if err != nil {
			if errors.Is(err, postgres_repository.ErrUserNotFound) || errors.Is(err, postgres_repository.ErrEmailOrPasswordWrong) {
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

func (h *AuthHandler) Logout(ctx *gin.Context) {
	go func() {
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user logouted!", map[string]interface{}{},
		)
	}()

	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
