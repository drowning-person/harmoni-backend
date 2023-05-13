package handler

import (
	followentity "harmoni/internal/entity/follow"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type FollowHandler struct {
	fs     *service.FollowService
	logger *zap.SugaredLogger
}

func NewFollowHandler(
	fs *service.FollowService,
	logger *zap.SugaredLogger,
) *FollowHandler {
	return &FollowHandler{
		fs:     fs,
		logger: logger,
	}
}

func (h *FollowHandler) Follow(c *fiber.Ctx) error {
	req := followentity.FollowRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.fs.Follow(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *FollowHandler) GetFollowers(c *fiber.Ctx) error {
	req := followentity.GetFollowersRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.fs.GetFollowers(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *FollowHandler) GetFollowings(c *fiber.Ctx) error {
	req := followentity.GetFollowingsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.fs.GetFollowing(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *FollowHandler) IsFollowing(c *fiber.Ctx) error {
	req := followentity.IsFollowingRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.fs.IsFollowing(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *FollowHandler) AreFollowEachOther(c *fiber.Ctx) error {
	req := followentity.AreFollowEachOtherRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.fs.AreFollowEachOther(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
