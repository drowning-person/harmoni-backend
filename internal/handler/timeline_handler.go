package handler

import (
	timelineentity "harmoni/internal/entity/timeline"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TimeLineHandler struct {
	ts *service.TimeLineService
}

func NewTimeLineHandler(ts *service.TimeLineService) *TimeLineHandler {
	return &TimeLineHandler{ts: ts}
}

func (h *TimeLineHandler) GetUserTimeLine(c *fiber.Ctx) error {
	req := timelineentity.GetUserTimeLineRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.ts.GetUserTimeLine(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}

func (h *TimeLineHandler) GetHomeTimeLine(c *fiber.Ctx) error {
	req := timelineentity.GetHomeTimeLineRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID
	reply, err := h.ts.GetHomeTimeLine(c.UserContext(), &req)

	return fiberx.HandleResponse(c, err, reply)
}
