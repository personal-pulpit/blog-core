package utils

type (
	SuccessfulHttpResponse struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
	}
	ErrorHttpResponse struct{
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Error    string `json:"error"`
	}
)
func NewSuccessfulHtppResponse(code int,msg string,data map[string]interface{})SuccessfulHttpResponse{
	return SuccessfulHttpResponse{Code: code,Message: msg,Data: data}
}
func NewErrorHtppResponse(code int,msg string,err error)ErrorHttpResponse{
	return ErrorHttpResponse{Code: code,Message: msg,Error: err.Error()}
}
