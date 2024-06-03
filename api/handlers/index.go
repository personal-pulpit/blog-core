package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/common"
	"blog/constants"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(ctx *gin.Context) {
	if common.GetUserStatus(ctx) {
		user, err := common.GetUserFromRedisById(helpers.GetIdFromToken(ctx))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, helpers.NewErrorHtppResponse(
				http.StatusInternalServerError, constants.MsgSometimeWentWrong, err,
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
