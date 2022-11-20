package model

import (
	"fiberLearn/pkg/zap"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	// 按照帖子时间排序
	OrderByTime = "time"
	// 按照点赞数量排序
	OrderByHot = "hot"
)

type ParamListData struct {
	PageSize int64  `query:"page_size"`
	PageNum  int64  `query:"page"`
	Order    string `query:"order" validate:"oneof=time,hot"`
}

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

func GetPageParam(c *fiber.Ctx) (*ParamListData, error) {
	param := &ParamListData{}
	if err := c.QueryParser(param); err != nil {
		zap.Logger.Error(err.Error())
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, err
		}
		return nil, errs
	}
	if param.PageNum < 0 {
		param.PageNum = 1
	} else if param.PageSize < 0 {
		param.PageSize = DefaultPageSize
	} else if param.PageSize > MaxPageSize {
		param.PageSize = MaxPageSize
	} else if param.Order != OrderByHot {
		param.Order = OrderByTime
	}
	return param, nil
}
