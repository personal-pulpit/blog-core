package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/common"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	if common.GetUserStatus(ctx) {
		ch := make(chan helpers.HttpResponse)
		go func() {
			user, err := common.GetUserFromRedisById(helpers.GetIdFromToken(ctx))
			if err != nil {
				ch <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), map[string]interface{}{},
				)
				return
			}
			ch <- helpers.NewHttpResponse(
				http.StatusOK, fmt.Sprintf("Hey %s", user["firstname"]), map[string]interface{}{},
			)
		}()
		helpers.GetResponse(ctx, http.StatusOK,ch)
	} else {
		ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
			http.StatusOK, "welcome to my api", map[string]interface{}{},
		))
	}
}
