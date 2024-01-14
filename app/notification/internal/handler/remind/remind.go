package remind

import (
	v1 "harmoni/app/notification/api/http/v1/notification"
	"harmoni/app/notification/internal/pkg/middleware"
	"harmoni/app/notification/internal/pkg/response"
	"harmoni/app/notification/internal/service/remind"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/action"
	"harmoni/internal/types/consts"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

type RemindHandler struct {
	rs     *remind.RemindService
	logger *log.Helper
}

func NewRemindHandler(
	rs *remind.RemindService,
	logger log.Logger,
) *RemindHandler {
	return &RemindHandler{
		rs:     rs,
		logger: log.NewHelper(log.With(logger, "module", "handler/notification")),
	}
}

func (h *RemindHandler) UnReadCount(c *gin.Context) {
	req := v1.UnReadRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.HandleResponse(c,
			errorx.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	req.UserID = middleware.GetUserInfoFromContext(c).GetId()
	resp, err := h.rs.UnreadCount(c, &req)
	response.HandleResponse(c, err, resp)
}

func (h *RemindHandler) ListLikeRemind(c *gin.Context) {
	req := v1.ListRemindRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.HandleResponse(c,
			errorx.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	req.Action = action.ActionLike
	req.UserID = middleware.GetUserInfoFromContext(c).GetId()
	resp, err := h.rs.ListRemind(c, &req)
	response.HandleResponse(c, err, resp)
}

func (h *RemindHandler) ListAtRemind(c *gin.Context) {
	req := v1.ListRemindRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.HandleResponse(c,
			errorx.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	req.Action = action.ActionAt
	req.UserID = middleware.GetUserInfoFromContext(c).GetId()
	resp, err := h.rs.ListRemind(c, &req)
	response.HandleResponse(c, err, resp)
}

func (h *RemindHandler) ListReplyRemind(c *gin.Context) {
	req := v1.ListRemindRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.HandleResponse(c,
			errorx.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	req.Action = action.ActionReply
	req.UserID = middleware.GetUserInfoFromContext(c).GetId()
	resp, err := h.rs.ListRemind(c, &req)
	response.HandleResponse(c, err, resp)
}

func (h *RemindHandler) ListLikeRemindDetail(c *gin.Context) {
	req := v1.LikeRemindDetailRequest{}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.HandleResponse(c,
			errorx.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	req.RemindID, err = strconv.ParseInt(c.Param("id"), consts.BaseDecimal, consts.BitSize64)
	if err != nil {
		response.HandleResponse(c,
			errorx.BadRequest(reason.RequestFormatError).WithError(err), nil)
		return
	}
	req.UserID = middleware.GetUserInfoFromContext(c).GetId()
	req.Action = action.ActionLike
	resp, err := h.rs.LikeDetail(c, &req)
	response.HandleResponse(c, err, resp)
}
