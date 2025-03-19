package helpers

import "github.com/gin-gonic/gin"

type (
	HttpResponse struct {
		Code    int                    `json:"code"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
)

func NewHttpResponse(code int, msg string, data map[string]interface{}) HttpResponse {
	return HttpResponse{Code: code, Message: msg, Data: data}
}

func RespondWithError(ctx *gin.Context, statusCode int, message string, metadata map[string]interface{}) {
	ctx.AbortWithStatusJSON(statusCode, NewHttpResponse(statusCode, message, metadata))
}

func RespondWithSuccess(ctx *gin.Context, statusCode int, message string, data map[string]interface{}) {
	ctx.JSON(statusCode, NewHttpResponse(statusCode, message, data))
}
