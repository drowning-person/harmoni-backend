package httpx

import "harmoni/internal/pkg/errorx"

type BaseBody struct {
	// http code
	Code int `json:"code"`
	// reason key
	Reason string `json:"reason"`
	// response message
	Message string `json:"msg"`
}

func NewBaseBody(code int, reason string) *BaseBody {
	return &BaseBody{
		Code:   code,
		Reason: reason,
	}
}

func NewBaseBodyFromError(e *errorx.Error) *BaseBody {
	return &BaseBody{
		Code:    int(e.Code),
		Reason:  e.Reason,
		Message: e.Message,
	}
}

// RespBody response body.
type RespBody struct {
	*BaseBody
	// response data
	Data interface{} `json:"data"`
}

// NewRespBody new response body
func NewRespBody(code int, reason string) *RespBody {
	return &RespBody{
		BaseBody: NewBaseBody(code, reason),
	}
}

// NewRespBodyFromError new response body from error
func NewRespBodyFromError(e *errorx.Error) *RespBody {
	return &RespBody{
		BaseBody: NewBaseBodyFromError(e),
	}
}

// NewRespBodyData new response body with data
func NewRespBodyData(code int, reason string, data interface{}) *RespBody {
	resp := NewRespBody(code, reason)
	resp.Data = data
	return resp
}

type PageResp struct {
	Total int64 `json:"total"`
	Page
}

type PageRespBody struct {
	*BaseBody
	Page PageResp    `json:"page"`
	List interface{} `json:"list"`
}

func NewPageRespBodyFromError(e *errorx.Error) *PageRespBody {
	return &PageRespBody{
		BaseBody: NewBaseBodyFromError(e),
	}
}

func NewPageRespBodyList(code int, reason string, total int64, page Page, list interface{}) *PageRespBody {
	return &PageRespBody{
		BaseBody: NewBaseBody(code, reason),
		Page: PageResp{
			Total: total,
			Page:  page,
		},
		List: list,
	}
}
