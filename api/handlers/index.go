package handlers

import (
	"blog/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(ctx *gin.Context) {
	if ctx.GetBool("is_logged") {
		claims := utils.GetToken(ctx)
		username := claims["username"]
		msg := fmt.Sprintf("Hey %s", username)
		ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
			http.StatusOK, msg, map[string]interface{}{},
		))
		return
	}
	ctx.JSON(http.StatusOK, utils.NewSuccessfulHtppResponse(
		http.StatusOK, "welcome to my api",map[string]interface{}{},
	))

}
