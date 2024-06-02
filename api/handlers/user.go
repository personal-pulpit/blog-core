package handlers

import (
	"blog/api/helpers"
	"blog/pkg/data/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		UserRepo *repo.UserRepo
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

func (u *User) GetAll(ctx *gin.Context) {
	users, err := u.UserRepo.GetAll()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting users!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "users got!", map[string]interface{}{
			"users": users,
		},
	))
}

func (u *User) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := u.UserRepo.GetById(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in getting user!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "user Got!", map[string]interface{}{
			"firstname":    user["firstname"],
			"lastname":     user["lastname"],
			"biography":    user["biography"],
			"username":     user["username"],
			"email":        user["email"],
			"phone number": user["phonenumber"],
			"created at":   user["createdAt"],
			"updated at":   user["updatedAt"],
		},
	))
}
func (u *User) Verify(ctx *gin.Context) {
	var li LoginInput
	err := ctx.ShouldBind(&li)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "sometimes went wrong", err),
		)
		return
	}
	user, err := u.UserRepo.Verify(li.Username, li.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "password or username is wrong", err),
		)
		return
	}
	err = helpers.SetToken(ctx, user.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed...", err),
		)
		return
	}
	ctx.Set("is_logged", true)
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "welcome back!", map[string]interface{}{
			"firstname":    user.Firstname,
			"lastname":     user.Lastname,
			"biography":    user.Biography,
			"username":     user.Username,
			"email":        user.Email,
			"phone number": user.PhoneNumber,
		},
	))

}
func (u *User) Create(ctx *gin.Context) {
	var si SigninInput
	err := ctx.ShouldBind(&si)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "sometimes went wrong", err),
		)
		return
	}
	user, err := u.UserRepo.Create(
		si.Firstname,
		si.Lastname,
		si.Biography,
		si.Username,
		si.Password,
		si.Email,
		si.PhoneNumber,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in creatig user", err),
		)
		return
	}
	err = helpers.SetToken(ctx, user.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed...", err),
		)
		return
	}
	ctx.Set("is_logged", true)
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
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
	id := helpers.GetIdFromToken(ctx)
	var ui UpdateInput
	err := ctx.ShouldBind(&ui)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "sometimes went wrong", err),
		)
		return
	}
	user, err := u.UserRepo.UpdateById(
		id,
		ui.Firstname,
		ui.Lastname,
		ui.Biography,
		ui.Username,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in updating user!", err),
		)
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
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
	err := u.UserRepo.Delete(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.NewErrorHtppResponse(
			http.StatusBadRequest, "failed in deleting user!", err))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "user deleted!", map[string]interface{}{},
	))
}
func (u *User) Logout(ctx *gin.Context) {
	u.UserRepo.DeleteChacheById(helpers.GetIdFromToken(ctx))
	helpers.DestroyToken(ctx)
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "user logouted!", map[string]interface{}{},
	))
}
