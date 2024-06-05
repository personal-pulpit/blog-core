package helpers

type (
	HttpResponse struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
	}

)
func NewHttpResponse(code int,msg string,data map[string]interface{})HttpResponse{
	return HttpResponse{Code: code,Message: msg,Data: data}
}
