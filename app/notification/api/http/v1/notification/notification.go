package v1

import "harmoni/internal/types/action"

type UnReadRequest struct {
	UserID int64         `json:"user_id"`
	Action action.Action `form:"action"`
}

type UnReadResponse struct {
	Count int64 `json:"count"`
}
