package apis

import (
	"fiberLearn/model"
	"fiberLearn/pkg/app"
	"fiberLearn/pkg/errcode"
	"fiberLearn/pkg/validator"
	"fiberLearn/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetTags(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	param, err := model.GetPageParam(c)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, total, err1 := services.GetTags(param)
	if err != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponseList(data, total)
}

func CreateTag(c *fiber.Ctx) error {
	var service services.TagInsertService
	r := app.NewResponse(c)
	if err := c.BodyParser(&service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams)
	}
	if err := validator.Validate(service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, err := service.Insert()
	if err != nil {
		return r.ToErrorResponse(err)
	}
	return r.ToResponse(data)
}

func GetTagDetail(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	id := c.Params("id")
	tagID, err := strconv.Atoi(id)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams)
	}

	data, err1 := services.GetTagDetail(tagID)
	if err1 != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponse(data)
}
