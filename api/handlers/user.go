package handlers

import (
	"blog/api/helpers"
	mysql_repository "blog/database/mysql/repo"

	"blog/internal/repository"

	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		UserMysqlRepo repository.UserMysqlRepository
		UserRedisRepo repository.UserRedisRepository
	}
	UpdateInput struct {
		Firstname string `form:"firstname" binding:"required"`
		Lastname  string `form:"lastname" binding:"required"`
		Username  string `form:"username" binding:"required,usernamevalidaitor"`
		Biography string `form:"biography" binding:"required"`
	}
	SigninInput struct {
		Firstname   string `form:"firstname" binding:"required"`
		Lastname    string `form:"lastname" binding:"required"`
		Username    string `form:"username" binding:"required,usernamevalidaitor"`
		Password    string `form:"password" binding:"required"`
		Email       string `form:"email" binding:"required,emailvalidatior"`
		PhoneNumber string `form:"phonenumber" binding:"required,phonenumbervalidaitor"`
		Biography   string `form:"biography" binding:"required"`
	}
	LoginInput struct {
		Username string `form:"username" binding:"required,usernamevalidaitor"`
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

func (u *User) GetAll(ctx *gin.Context) {
	go func() {
		users, err := u.UserRedisRepo.GetCaches()
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(http.StatusOK, "users got!", map[string]interface{}{"users": users})
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) GetByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		user, err := u.UserRedisRepo.GetCacheByID(ID)
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
				"firstname":    user["firstname"],
				"lastname":     user["lastname"],
				"biography":    user["biography"],
				"username":     user["username"],
				"email":        user["email"],
				"phone number": user["phonenumber"],
				"created at":   user["createdAt"],
				"updated at":   user["updatedAt"],
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
		user, err := u.UserMysqlRepo.Verify(li.Username, li.Password)
		if err != nil {
			if errors.Is(err, mysql_repository.ErrUserNotFound) || errors.Is(err, mysql_repository.ErrUsernameOrPasswordWrong) {
				userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		err = helpers.SetToken(ctx, user.ID)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(http.StatusInternalServerError, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "welcome back!", map[string]interface{}{
				"firstname":    user.Firstname,
				"lastname":     user.Lastname,
				"biography":    user.Biography,
				"username":     user.Username,
				"email":        user.Email,
				"phone number": user.PhoneNumber,
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
			} else if utils.CheckErrorForWord(err, "usernamevalidaitor") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrUsernameShouldContain),
					nil)
				return
			} else if utils.CheckErrorForWord(err, "emailvalidatior") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrInvalidEmail),
					nil)
				return
			} else if utils.CheckErrorForWord(err, "phonenumbervalidaitor") {
				userResponseChannel <- helpers.NewHttpResponse(
					http.StatusBadRequest,
					utils.GetValidationError(ErrInvalidPhonenumber),
					nil)
				return
			}
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest,
				utils.GetValidationError(err),
				nil)
			return
		}
		user, tx, err := u.UserMysqlRepo.Create(
			si.Firstname,
			si.Lastname,
			si.Biography,
			si.Username,
			si.Password,
			si.Email,
			si.PhoneNumber,
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
		err = helpers.SetToken(ctx, user.ID)
		if err != nil {
			tx.Rollback()
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusInternalServerError, err.Error(), nil)
			return
		}
		tx.Commit()
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusCreated, "user created!", map[string]interface{}{
				"firstname":    user.Firstname,
				"lastname":     user.Lastname,
				"biography":    user.Biography,
				"username":     user.Username,
				"email":        user.Email,
				"phone number": user.PhoneNumber,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusCreated, userResponseChannel)
}
func (u *User) UpdateByID(ctx *gin.Context) {
	go func() {
		ID := helpers.GetIdFromToken(ctx)
		var ui UpdateInput
		err := ctx.ShouldBind(&ui)
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
		user, err := u.UserMysqlRepo.UpdateByID(
			ID,
			ui.Firstname,
			ui.Lastname,
			ui.Biography,
			ui.Username,
		)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user updated!", map[string]interface{}{
				"firstname": user.Firstname,
				"lastname":  user.Lastname,
				"biography": user.Biography,
				"username":  user.Username,
			},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) DeleteByID(ctx *gin.Context) {
	go func() {
		ID := ctx.Param("ID")
		err := u.UserMysqlRepo.DeleteByID(ID)
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
		u.UserRedisRepo.DeleteCacheByID(helpers.GetIdFromToken(ctx))
		helpers.DestroyToken(ctx)
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user logouted!", map[string]interface{}{},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
