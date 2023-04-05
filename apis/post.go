package apis

import (
	"harmoni/model"
	"harmoni/pkg/app"
	"harmoni/pkg/errcode"
	"harmoni/pkg/validator"
	"harmoni/pkg/zap"
	"harmoni/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetPostDetail(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	id := c.Params("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams)
	}

	data, err1 := services.GetPostDetail(int64(postID))
	if err1 != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponse(data)
}

func GetPosts(c *fiber.Ctx) error {
	r := app.NewResponse(c)
	param, err := model.GetPageParam(c)
	if err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	data, total, err1 := services.GetPosts(param)
	if err != nil {
		return r.ToErrorResponse(err1)
	}
	return r.ToResponseList(data, total)
}

func CreatePost(c *fiber.Ctx) error {
	var service services.PostInsertService
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

func LikePost(c *fiber.Ctx) error {
	var service services.PostLikeService
	r := app.NewResponse(c)
	if err := c.BodyParser(&service); err != nil {
		zap.Logger.Error(err.Error())
		return r.ToErrorResponse(errcode.InvalidParams)
	}
	if err := validator.Validate(service); err != nil {
		return r.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
	}
	if err := service.LikePost(int64(c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64))); err != nil {
		return r.ToErrorResponse(err)
	}
	return r.ToResponse(nil)
}
