package handler

import (
	likeentity "harmoni/app/harmoni/internal/entity/like"
	"harmoni/app/harmoni/internal/pkg/fiberx"
	"harmoni/app/harmoni/internal/pkg/middleware"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/app/harmoni/internal/service"
	"harmoni/internal/pkg/errorx"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LikeHandler struct {
	ls     *service.LikeService
	logger *zap.SugaredLogger
}

func NewLikeHandler(
	ls *service.LikeService,
	logger *zap.SugaredLogger,
) *LikeHandler {
	return &LikeHandler{
		ls:     ls,
		logger: logger.With("module", "handler/like"),
	}
}

func (h *LikeHandler) Like(c *fiber.Ctx) error {
	req := likeentity.LikeRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		h.logger.Warn(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.ls.Like(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *LikeHandler) IsLiking(c *fiber.Ctx) error {
	req := likeentity.IsLikingRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		h.logger.Warn(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.ls.IsLiking(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *LikeHandler) LikingList(c *fiber.Ctx) error {
	req := likeentity.GetLikingsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		h.logger.Warn(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.ls.GetLikings(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
