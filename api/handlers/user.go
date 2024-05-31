package handlers

import (
	"blog/pkg/data/repo"
	"blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserRepo *repo.UserRepo
}

func (u *User) GetUsers(ctx *gin.Context) {
	users, err := u.UserRepo.GetUsersRedis()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting users!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "users got!", map[string]interface{}{
			"user": users,
		},
	))
}

func (u *User) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := u.UserRepo.GetUserByIdRedis(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting user!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user Got!", map[string]interface{}{
			"firstname":    user["firstname"],
			"lastname":     user["lastname"],
			"biography":    user["biography"],
			"username":     user["username"],
			"email":        user["email"],
			"phone number": user["phonenumber"],
			"created at": user["createdAt"],
			"updated at": user["updatedAt"],
		},
	))
}
func (u *User) Verify(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	user, err := u.UserRepo.VerifyUser(username, password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "password or username is wrong", err),
		)
		return
	}
	err = utils.SetToken(ctx, user.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed...", err),
		)
		return
	}
	ctx.Set("is_logged", true)
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "welcome back!", map[string]interface{}{
			"firstname":    user.Firstname,
			"lastname":     user.Lastname,
			"biography":user.Biography,
			"username":     user.Username,
			"email":        user.Email,
			"phone number": user.PhoneNumber,
		},
	))

}
func (u *User) Create(ctx *gin.Context) {
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	email := ctx.PostForm("email")
	phonenumber := ctx.PostForm("phonenumber")
	biography := ctx.PostForm("biography")
	user, err := u.UserRepo.CreateUser(firstname, lastname, biography, username, password, email, phonenumber)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig user", err),
		)
		return
	}
	err = utils.SetToken(ctx, user.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed...", err),
		)
		return
	}
	ctx.Set("is_logged", true)
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user created!", map[string]interface{}{
			"firstname":    user.Firstname,
			"lastname":     user.Lastname,
			"biography":    user.Biography,
			"username":     user.Username,
			"email":        user.Email,
			"phone number": user.PhoneNumber,
		},
	))
}
func (u *User) UpdateById(ctx *gin.Context) {
	id := ctx.PostForm("id")
	firstname := ctx.PostForm("firstname")
	lastname := ctx.PostForm("lastname")
	username := ctx.PostForm("username")
	biography := ctx.PostForm("biography")
	user, err := u.UserRepo.UpdateUserById(id, firstname, lastname, biography, username)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in updating user!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user updated!", map[string]interface{}{
			"firstname": user.Firstname,
			"lastname":  user.Lastname,
			"biography": user.Biography,
			"username":  user.Username,
		},
	))
}
func (u *User) DeleteById(ctx *gin.Context) {
	id := ctx.Param("id")

	err := u.UserRepo.DeleteUser(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in deleting user!", err))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user deleted!", map[string]interface{}{},
	))
}
func (u *User) Logout(ctx *gin.Context) {
	u.UserRepo.DeleteChacheByIdRedis(utils.GetIdFromToken(ctx))
	utils.DestroyToken(ctx)
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "user logouted!", map[string]interface{}{},
	))
}
