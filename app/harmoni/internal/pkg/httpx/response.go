package httpx

import "harmoni/app/harmoni/internal/pkg/errorx"

// RespBody response body.
type RespBody struct {
	// http code
	Code int `json:"code"`
	// reason key
	Reason string `json:"reason"`
	// response message
	Message string `json:"msg"`
	// response data
	Data interface{} `json:"data"`
}

// NewRespBody new response body
func NewRespBody(code int, reason string) *RespBody {
	return &RespBody{
		Code:   code,
		Reason: reason,
	}
}

// NewRespBodyFromError new response body from error
func NewRespBodyFromError(e *errorx.Error) *RespBody {
	return &RespBody{
		Code:    e.Code,
		Reason:  e.Reason,
		Message: e.Message,
	}
}

// NewRespBodyData new response body with data
func NewRespBodyData(code int, reason string, data interface{}) *RespBody {
	return &RespBody{
		Code:   code,
		Reason: reason,
		Data:   data,
	}
}
