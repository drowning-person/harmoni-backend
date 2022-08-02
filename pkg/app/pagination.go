package app

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	DefaultPageSize = 25
	MaxPageSize     = 100
)

func GetPage(c *fiber.Ctx) int {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *fiber.Ctx) int {
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}

	return pageSize
}

func GetPageOffset(c *fiber.Ctx) (offset, limit int) {
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ = strconv.Atoi(c.Query("page_size"))
	if limit <= 0 {
		limit = DefaultPageSize
	} else if limit > MaxPageSize {
		limit = MaxPageSize
	}
	offset = (page - 1) * limit
	return
}
