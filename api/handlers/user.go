package handlers

import (
	"blog/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct{}

func (u User) GetAll(ctx *gin.Context) {}

func (u User) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := service.GetUser(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (u User) Create(ctx *gin.Context) {
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")
	phonenumber := ctx.PostForm("phonenumber")
	user, err := service.CreateUser(firstname, lastname, username, password, email, phonenumber)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}
func (u User) UpdateById(ctx *gin.Context) {
	id := ctx.PostForm("id")
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")
	phonenumber := ctx.PostForm("phonenumber")

	err := service.UpdateUserById(id, firstname, lastname, username, password, email, phonenumber)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "ok!"})
}
func (u User) DeleteById(ctx *gin.Context) {
	id := ctx.Param("id")

	err := service.DeleteUser(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": "ok!"})
}
