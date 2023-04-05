package apis

import (
	"harmoni/model"
	"harmoni/pkg/app"
	"harmoni/pkg/errcode"
	"harmoni/pkg/validator"
	"harmoni/pkg/zap"
	"harmoni/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CreateComment(c *fiber.Ctx) error {
	var service services.CommentInsertService
	r := app.NewResponse(c)
	if err := c.BodyParser(&service); err != nil {
		zap.Logger.Error(err.Error())
		return r.ToErrorResponse(errcode.InvalidParams)
	}
	if err := validator.Validate(service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, err := service.Insert(int64(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)))
	if err != nil {
		return r.ToErrorResponse(err)
	}
	return r.ToResponse(data)
}

func GetPostComment(c *fiber.Ctx) error {
	var service services.GetPostCommentsService
	r := app.NewResponse(c)
	if err := c.QueryParser(&service); err != nil {
		zap.Logger.Error(err.Error())
		return r.ToErrorResponse(errcode.InvalidParams)
	}
	param, err := model.GetPageParam(c)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	if err := validator.Validate(service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, total, err1 := service.Retrieve(*param)
	if err != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponseList(data, total)
}
