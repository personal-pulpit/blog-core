package handlers

import (
	"blog/api/helpers"
	"blog/pkg/data/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	User struct {
		UserRepo repo.UserDB
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

var userResponseChannel = make(chan helpers.HttpResponse)

func (u *User) GetAll(ctx *gin.Context) {
	go func() {
		users, err := u.UserRepo.GetAll()
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
			return
		}
		userResponseChannel <- helpers.NewHttpResponse(http.StatusOK, "users got!", map[string]interface{}{"users": users})
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
func (u *User) GetById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		user, err := u.UserRepo.GetById(id)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
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
			userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, err := u.UserRepo.Verify(li.Username, li.Password)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
			return
		}
		err = helpers.SetToken(ctx, user.Id)
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
			userResponseChannel <- helpers.NewHttpResponse(http.StatusBadRequest, err.Error(), nil)
			return
		}
		user, tx, err := u.UserRepo.Create(
			si.Firstname,
			si.Lastname,
			si.Biography,
			si.Username,
			si.Password,
			si.Email,
			si.PhoneNumber,
		)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
			return
		}
		err = helpers.SetToken(ctx, user.Id)
		if err != nil {
			tx.Rollback()
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
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
func (u *User) UpdateById(ctx *gin.Context) {
	go func() {
		id := helpers.GetIdFromToken(ctx)
		var ui UpdateInput
		err := ctx.ShouldBind(&ui)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
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
func (u *User) DeleteById(ctx *gin.Context) {
	go func() {
		id := ctx.Param("id")
		err := u.UserRepo.DeleteById(id)
		if err != nil {
			userResponseChannel <- helpers.NewHttpResponse(
				http.StatusBadRequest, err.Error(), nil)
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
		u.UserRepo.DeleteChacheById(helpers.GetIdFromToken(ctx))
		helpers.DestroyToken(ctx)
		userResponseChannel <- helpers.NewHttpResponse(
			http.StatusOK, "user logouted!", map[string]interface{}{},
		)
	}()
	helpers.GetResponse(ctx, http.StatusOK, userResponseChannel)
}
