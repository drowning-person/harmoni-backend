package app

import (
	"fiberLearn/pkg/errcode"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Ctx *fiber.Ctx
}

type Pager struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

func NewResponse(ctx *fiber.Ctx) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = fiber.Map{
			"code": 0,
			"msg":  "success",
		}
	} else {
		data = fiber.Map{
			"code": 0,
			"msg":  "success",
			"data": data,
		}
	}
	r.Ctx.Status(http.StatusOK).JSON(data)
}

func (r *Response) ToResponseList(list interface{}, total int64) {
	r.ToResponse(fiber.Map{
		"list": list,
		"pager": Pager{
			Page:     GetPage(r.Ctx),
			PageSize: GetPageSize(r.Ctx),
			Total:    total,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := fiber.Map{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.Status(err.StatusCode()).JSON(response)
}
