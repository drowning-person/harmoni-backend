package apis

import (
	"fiberLearn/pkg/app"
	"fiberLearn/pkg/errcode"
	"fiberLearn/pkg/validator"
	"fiberLearn/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetTags(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	offset, limit := app.GetPageOffset(c)
	data, total, err := services.GetTags(offset, limit)
	if err != nil {
		return r.ToErrorResponse(errcode.GetTagsFailed)
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
		return r.ToErrorResponse(errcode.CreateTagFailed)
	}
	if data == nil {
		return r.ToErrorResponse(errcode.TagHasExisted)
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

	data, err := services.GetTagDetail(tagID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.ToErrorResponse(errcode.NotFound.WithDetails("话题不存在"))
		}
		return r.ToErrorResponse(errcode.GetTagFailed)
	}
	return r.ToResponse(data)
}
