package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Base(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"data": "ok!"})
}
