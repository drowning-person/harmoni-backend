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

func GetAllUsers(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	offset, limit := app.GetPageOffset(c)
	data, total, err := services.GetUsers(offset, limit)
	if err != nil {
		return r.ToErrorResponse(errcode.GetUsersFailed)
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

	data, err := services.GetUser(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.ToErrorResponse(errcode.NotFound.WithDetails("用户不存在"))
		}
		return r.ToErrorResponse(errcode.GetUserFailed)
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
	data, err := service.Regist()
	if err != nil {
		return r.ToErrorResponse(errcode.UserRegisterFailed)
	}
	if data == nil {
		return r.ToErrorResponse(errcode.UsernameHasExisted)
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
	data, token, err := service.Login()
	if err != nil {
		return r.ToErrorResponse(err.(*errcode.Error))
	}
	return r.ToResponse((fiber.Map{"token": token, "user": data}))
}
