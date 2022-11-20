package app

import (
	"fiberLearn/model"
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

func (r *Response) ToResponse(data interface{}) error {
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
	return r.Ctx.Status(http.StatusOK).JSON(data)
}

func (r *Response) ToResponseList(list interface{}, total int64) error {
	return r.ToResponse(fiber.Map{
		"list": list,
		"pager": Pager{
			Page:     model.GetPage(r.Ctx),
			PageSize: model.GetPageSize(r.Ctx),
			Total:    total,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) error {
	response := fiber.Map{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	return r.Ctx.Status(err.StatusCode()).JSON(response)
}
