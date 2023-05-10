package handler

import (
	accountentity "harmoni/internal/entity/account"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AccountHandler struct {
	as     *service.AccountService
	jwtmw  *middleware.JwtAuthMiddleware
	logger *zap.SugaredLogger
}

func NewAccountHandler(as *service.AccountService, jwtmw *middleware.JwtAuthMiddleware, logger *zap.SugaredLogger) *AccountHandler {
	return &AccountHandler{
		as:     as,
		jwtmw:  jwtmw,
		logger: logger,
	}
}

func (h *AccountHandler) MailSend(c *fiber.Ctx) error {
	req := accountentity.MailSendRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	if req.Type != accountentity.RegisterAct {
		if err := h.jwtmw.ParseAndVerifyToken(c); err != nil {
			return fiberx.HandleResponse(c, err, nil)
		}
		req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	}
	reply, err := h.as.MailSend(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *AccountHandler) MailCheck(c *fiber.Ctx) error {
	req := accountentity.MailCheckRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	if req.Type != accountentity.RegisterAct {
		if err := h.jwtmw.ParseAndVerifyToken(c); err != nil {
			return fiberx.HandleResponse(c, err, nil)
		}
		req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	}
	reply, err := h.as.MailCheck(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *AccountHandler) ChangeEmail(c *fiber.Ctx) error {
	req := accountentity.ChangeEmailRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.as.ChangeEmail(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *AccountHandler) ChangePassword(c *fiber.Ctx) error {
	req := accountentity.ChangePasswordRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.as.ChangePassword(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *AccountHandler) ResetPassword(c *fiber.Ctx) error {
	req := accountentity.ResetPasswordRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.as.ResetPassword(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
