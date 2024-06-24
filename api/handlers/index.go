package handlers

import (
	"blog/api/helpers"
	"blog/api/helpers/auth_helper"
	"blog/api/helpers/common"
	"blog/pkg/auth_manager"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Main struct {
	AuthHelper auth_helper.AuthHeaderHelper
}

func (h *Main) Main(ctx *gin.Context) {
	if common.GetUserStatus(ctx) {
		ch := make(chan helpers.HttpResponse)
		go func() {
			token, _ := h.AuthHelper.GetHeader(ctx)
			id,err := auth_helper.GetIdByToken(token, auth_manager.RefreshToken)
			if err != nil {
				ch <- helpers.NewHttpResponse(
					http.StatusInternalServerError, err.Error(), map[string]interface{}{},
				)
				return
			}
			user, err := common.GetUserFromRedisByID(id)
			if err != nil {
				ch <- helpers.NewHttpResponse(
					http.StatusBadRequest, err.Error(), map[string]interface{}{},
				)
				return
			}
			ch <- helpers.NewHttpResponse(
				http.StatusOK, fmt.Sprintf("Hey %s", user["FirstName"]), map[string]interface{}{},
			)
		}()
		helpers.GetResponse(ctx, http.StatusOK, ch)
	} else {
		ctx.JSON(http.StatusOK, helpers.NewHttpResponse(
			http.StatusOK, "welcome to my api", map[string]interface{}{},
		))
	}
}
