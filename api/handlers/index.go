package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/common"
	"blog/internal/service/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Main struct{
	UserService user.UserService
}

func (h *Main) Main(ctx *gin.Context) {
	if common.GetUserStatus(ctx) {
		ch := make(chan helpers.HttpResponse)
		go func() {
			id := ctx.GetString("id")
			user, err := h.UserService.GetUserProfile(id)
			if err != nil {
				ch <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), map[string]interface{}{},
				)
				return
			}
			ch <- helpers.NewHttpResponse(
				http.StatusOK, fmt.Sprintf("Hey %s %s", user.FirstName,user.LastName), map[string]interface{}{},
			)
		}()
		helpers.GetResponse(ctx, http.StatusOK, ch)
	} else {
		ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
			http.StatusOK, "welcome to my api", map[string]interface{}{},
		))
	}
}
