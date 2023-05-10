package handler

import (
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	us *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		us: us,
	}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	req := userentity.GetUsersRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.us.GetUsers(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	req := userentity.GetUserDetailRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.us.GetUserByUserID(c.UserContext(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	req := userentity.UserRegisterRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.us.RegisterByEmail(c.UserContext(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	req := userentity.UserLoginRequset{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.us.Login(c.UserContext(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	req := userentity.UserLogoutRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	myclaims := middleware.GetClaimsFromCtx(c.UserContext())
	req.UserID = myclaims.UserID
	req.AccessTokenID = myclaims.ID

	reply, err := h.us.Logout(c.UserContext(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *UserHandler) RefreshToken(c *fiber.Ctx) error {
	req := userentity.RefreshTokenRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.us.RefreshToken(c.UserContext(), &req)
	return fiberx.HandleResponse(c, err, reply)
}
