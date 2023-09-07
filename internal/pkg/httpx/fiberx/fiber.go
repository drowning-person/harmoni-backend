package fiberx

import (
	"errors"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/validator"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// HandleResponse Handle response body
func HandleResponse(c *fiber.Ctx, err error, data interface{}) error {
	// no error
	if err == nil {
		c.Status(http.StatusOK).JSON(httpx.NewRespBodyData(http.StatusOK, reason.Success, data))
		return nil
	}

	var myErr *errorx.Error
	// unknown error
	if !errors.As(err, &myErr) {
		zap.L().Sugar().Error(err, "\n", errorx.LogStack(2, 5))
		c.Status(http.StatusInternalServerError).JSON(httpx.NewRespBody(
			http.StatusInternalServerError, reason.UnknownError))
		return err
	}

	respBody := httpx.NewRespBodyFromError(myErr)
	if data != nil {
		respBody.Data = data
	}

	c.Status(myErr.Code).JSON(respBody)
	return nil
}

func Parser(c *fiber.Ctx, out interface{}) error {
	method := c.Method()
	if method != http.MethodHead && method != http.MethodGet {
		if err := c.BodyParser(out); err != nil {
			if err != fiber.ErrUnprocessableEntity {
				return err
			}
		}
	}
	if err := c.ParamsParser(out); err != nil {
		return err
	}
	if err := c.QueryParser(out); err != nil {
		return err
	}
	return nil
}

func ParseAndCheck(c *fiber.Ctx, out interface{}) error {
	if err := Parser(c, out); err != nil {
		return err
	}
	if err := validator.Validate(out); err != nil {
		return err
	}

	return nil
}
