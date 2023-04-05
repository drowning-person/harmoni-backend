package apis

import (
	"harmoni/model"
	"harmoni/pkg/app"
	"harmoni/pkg/errcode"
	"harmoni/pkg/validator"
	"harmoni/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllUsers(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	param, err := model.GetPageParam(c)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, total, err1 := services.GetUsers(param)
	if err != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponseList(data, total)
}

func GetUser(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	id := c.Params("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams)
	}

	data, err1 := services.GetUser(userID)
	if err1 != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponse(data)
}

func Regist(c *fiber.Ctx) error {
	var service services.UserRegistService
	r := app.NewResponse(c)
	if err := c.BodyParser(&service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams)
	}
	if err := validator.Validate(service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, err1 := service.Regist()
	if err1 != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponse(data)
}

func Login(c *fiber.Ctx) error {
	var service services.UserLoginService
	r := app.NewResponse(c)
	if err := c.BodyParser(&service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams)
	}
	if err := validator.Validate(service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, token, err1 := service.Login()
	if err1 != nil {
		return r.ToErrorResponse(err1.(*errcode.Error))
	}
	return r.ToResponse((fiber.Map{"token": token, "user": data}))
}
