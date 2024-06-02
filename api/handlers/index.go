package handlers

import (
	"blog/api/helpers/common"
	"blog/api/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(ctx *gin.Context) {
	if common.GetUserStatus(ctx) {
		id := helpers.GetIdFromToken(ctx)
		user, err := common.GetUserFromRedisById(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helpers.NewErrorHtppResponse(
				http.StatusInternalServerError, "sometime went wrong!", err,
			))
			return
		}
		msg := fmt.Sprintf("Hey %s", user["firstname"])
		ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
			http.StatusOK, msg, map[string]interface{}{},
		))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewSuccessfulHtppResponse(
		http.StatusOK, "welcome to my api", map[string]interface{}{},
	))

}
