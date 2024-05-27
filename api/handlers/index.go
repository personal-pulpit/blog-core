package handlers

import (
	"blog/pkg/service"
	"blog/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(ctx *gin.Context) {
	if ctx.GetBool("is_logged") {
		id := utils.GetIdFromToken(ctx)
		user,err := service.GetUserByIdRedis(id)
		if err != nil{
			ctx.JSON(http.StatusInternalServerError, utils.NewErrorHtppResponse(
				http.StatusInternalServerError, "sometime went wrong!", err,
			))
			return
		}
		msg := fmt.Sprintf("Hey %s", user.Firstname)
		ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
			http.StatusOK, msg, map[string]interface{}{},
		))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "welcome to my api",map[string]interface{}{},
	))

}
