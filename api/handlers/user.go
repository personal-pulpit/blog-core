package handlers

import (
	"blog/pkg/service"
	"blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct{}

func (u User) GetUsers(ctx *gin.Context) {
	users, err := service.GetUsers()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,  utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting users!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user created!", map[string]interface{}{
			"user":users,
		},
	))
}

func (u User) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := service.GetUser(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,  utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting user!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user Got!", map[string]interface{}{
			"fristname":    user.Fristname,
			"lastname":     user.Lastname,
			"username":     user.Username,
			"email":        user.Email,
			"phone number": user.PhoneNumber,
		},
	))
}

func (u User) Create(ctx *gin.Context) {
	firstname := ctx.PostForm("fristname")
	lastname := ctx.PostForm("lastname")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")
	phonenumber := ctx.PostForm("phonenumber")
	user, err := service.CreateUser(firstname, lastname, username, password, email, phonenumber)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig user", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user created!", map[string]interface{}{
			"fristname":    user.Fristname,
			"lastname":     user.Lastname,
			"username":     user.Username,
			"email":        user.Email,
			"phone number": user.PhoneNumber,
		},
	))
}
func (u User) UpdateById(ctx *gin.Context) {
	id := ctx.PostForm("id")
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")
	phonenumber := ctx.PostForm("phonenumber")

	user,err := service.UpdateUserById(id, firstname, lastname, username, password, email, phonenumber)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,  utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in updating user!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user updated!", map[string]interface{}{
			"fristname":    user.Fristname,
			"lastname":     user.Lastname,
			"username":     user.Username,
			"email":        user.Email,
			"phone number": phonenumber,
		},
	))
}
func (u User) DeleteById(ctx *gin.Context) {
	id := ctx.Param("id")

	err := service.DeleteUser(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in deleting user!", err),)
			return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user updated!", map[string]interface{}{},
	))
}
