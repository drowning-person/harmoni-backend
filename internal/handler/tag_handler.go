package handler

import (
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TagHandler struct {
	ts *service.TagService
}

func NewTagHandler(ts *service.TagService) *TagHandler {
	return &TagHandler{
		ts: ts,
	}
}

func (h *TagHandler) GetTags(c *fiber.Ctx) error {
	req := tagentity.GetTagsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.ts.GetTags(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	req := tagentity.CreateTagRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.ts.Create(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *TagHandler) GetTagByID(c *fiber.Ctx) error {
	req := tagentity.GetTagDetailRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.ts.GetTag(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
