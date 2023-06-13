package handler

import (
	postentity "harmoni/internal/entity/post"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	ps *service.PostService
}

func NewPostHandler(ps *service.PostService) *PostHandler {
	return &PostHandler{ps: ps}
}

func (h *PostHandler) GetPostDetail(c *fiber.Ctx) error {
	req := postentity.GetPostDetailRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.ps.GetPostDetail(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *PostHandler) GetPosts(c *fiber.Ctx) error {
	req := postentity.GetPostsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
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
