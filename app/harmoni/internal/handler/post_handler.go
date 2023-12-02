package handler

import (
	postentity "harmoni/app/harmoni/internal/entity/post"
	"harmoni/app/harmoni/internal/pkg/fiberx"
	"harmoni/app/harmoni/internal/pkg/middleware"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/app/harmoni/internal/service"
	"harmoni/internal/pkg/errorx"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	ps *service.PostService
}

func NewPostHandler(ps *service.PostService) *PostHandler {
	return &PostHandler{ps: ps}
}

func (h *PostHandler) GetPostInfo(c *fiber.Ctx) error {
	req := postentity.GetPostInfoRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	claim := middleware.GetClaimsFromCtx(c.UserContext())
	if claim != nil {
		req.UserID = claim.UserID
	}
	reply, err := h.ps.GetPostInfo(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *PostHandler) GetPosts(c *fiber.Ctx) error {
	req := postentity.GetPostsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	claim := middleware.GetClaimsFromCtx(c.UserContext())
	if claim != nil {
		req.UserID = claim.UserID
	}
	reply, err := h.ps.GetPosts(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *PostHandler) CreatePost(c *fiber.Ctx) error {
	req := postentity.CreatePostRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.ps.Create(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
