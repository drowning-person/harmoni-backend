package handler

import (
	commententity "harmoni/app/harmoni/internal/entity/comment"
	"harmoni/app/harmoni/internal/pkg/fiberx"
	"harmoni/app/harmoni/internal/pkg/middleware"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/app/harmoni/internal/service"
	"harmoni/internal/pkg/errorx"

	"github.com/gofiber/fiber/v2"
)

type CommentHandler struct {
	cs *service.CommentService
}

func NewCommentHandler(cs *service.CommentService) *CommentHandler {
	return &CommentHandler{
		cs: cs,
	}
}

func (h *CommentHandler) GetComments(c *fiber.Ctx) error {
	req := commententity.GetCommentsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.cs.GetComments(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	req := commententity.CreateCommentRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.cs.Create(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
