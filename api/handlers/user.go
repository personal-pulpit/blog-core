package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	mysql_repository "blog/database/mysql/repo"
	"blog/pkg/auth_manager"

	"blog/internal/service/authentication"
	"blog/internal/service/user"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		AuthHelper auth_helper.AuthHeaderHelper
		AuthService authentication.AuthService
		UserService user.UserService
	}
	UpdateInput struct {
		FirstName string `form:"FirstName" binding:"required"`
		LastName  string `form:"lastname" binding:"required"`
		Biography string `form:"biography" binding:"required"`
	}
	SigninInput struct {
		FirstName   string `form:"FirstName" binding:"required"`
		LastName    string `form:"lastname" binding:"required"`
		Username    string `form:"username" binding:"required,usernamevalidaitor"`
		Password    string `form:"password" binding:"required"`
		Email       string `form:"email" binding:"required,emailvalidatior"`
		Biography   string `form:"biography" binding:"required"`
	}
	LoginInput struct {
		Email string `form:"email" binding:"required,emailvalidatior"`
		Password string `form:"password" binding:"required"`
	}
	DeleteAccountInput struct{
		Password string `form:"password" binding:"required"`
	}
)

var (
	userResponseChannel        = make(chan helpers.HttpResponse)
	ErrPleaseCompleteAllFields = errors.New("please complete all fields")
	ErrUsernameShouldContain   = errors.New("username should contain: a-z  _ 0-9")
	ErrInvalidEmail            = errors.New("email is invalid")
	ErrInvalidPhonenumber      = errors.New("phonenumber is invalid")
)
func (u *User) GetByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		user, err := u.UserService.GetUserProfile(ID)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrUserNotFound) {
				userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusCreated, "user Got!", map[string]interface{}{
				"user":user,
			},
		)

	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) Verify(ctx *gin.Context) {
	go func() {
		var li LoginInput
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
		user,accessToken,refreshToken,err := u.AuthService.Login(li.Email, li.Password)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrUserNotFound) || errors.Is(err, mysql_repository.ErrUsernameOrPasswordWrong) {
				userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "welcome back!", map[string]interface{}{
				"user":user,
				"access_token":accessToken,
				"refresh_token":refreshToken,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) Create(ctx *gin.Context) {
	go func() {
		var si SigninInput
		err := ctx.ShouldBind(&si)
		if err != nil {
			if utils.CheckErrorForWord(err, "required") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrPleaseCompleteAllFields),
					nil)
				return
			}  else if utils.CheckErrorForWord(err, "emailvalidatior") {
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
			si.FirstName,
			si.LastName,
			si.Email,
			si.Biography,
			si.Password,
		)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrEmailAlreadyExits) ||
				errors.Is(err, mysql_repository.ErrUsernameAlreadyExits) ||
				errors.Is(err, mysql_repository.ErrPhoneNumberAlreadyExits) {
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
				"user":user,
				"verifyEmailToken":verifyEmailToken,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusCreated, userResponseChannel)
}
func (u *User) UpdateByID(ctx *gin.Context) {
	go func() {
		token,err := u.AuthHelper.GetHeader(ctx)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		id,err := auth_helper.GetIdByToken(token,auth_manager.AccessToken)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		var ui UpdateInput
		err = ctx.ShouldBind(&ui)
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
		user,err := u.UserService.UpdateProfile(
			id,
			ui.FirstName,
			ui.LastName,
			ui.Biography,
		)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user updated!", map[string]interface{}{
				"user":user,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) DeleteByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		var si SigninInput
		err := ctx.ShouldBind(&si)
		err = u.UserService.DeleteAccount(ID,si.Password)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrUserNotFound) {
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
		token ,err := u.AuthHelper.GetHeader(ctx)
		if err != nil{
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
