package apis

import (
	"fiberLearn/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetTags(c *fiber.Ctx) error {
	data, err := services.GetTags()
	if err != nil {
		c.JSON(err)
		return err
	}
	return c.JSON(data)
}

func CreateTag(c *fiber.Ctx) error {
	var service services.TagInsertService

	if err := c.BodyParser(&service); err != nil {
		return err
	}
	data, err := service.Insert()
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func GetTagDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	tagID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(err)
		return err
	}
	service := services.TagDetailService{
		TagID: int64(tagID),
	}

	if err := service.Get(); err != nil {
		c.JSON(err)
		return err
	}
	return c.JSON(service)
}
