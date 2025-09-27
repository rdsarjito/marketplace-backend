package response

type BaseResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, err interface{}) BaseResponse {
	return BaseResponse{
		Status:  false,
		Message: message,
		Error:   err,
	}
}
