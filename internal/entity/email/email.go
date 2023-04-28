package email

import (
	"context"
	"encoding/json"
	"harmoni/internal/pkg/common"
	"time"
)

type EmailCodeContent struct {
	Code        string `json:"code,omitempty"`
	LastReqTime int64  `json:"last_req_time,omitempty"`
}

func (r *EmailCodeContent) ToJSONString() string {
	codeBytes, _ := json.Marshal(r)
	return string(codeBytes)
}

func (r *EmailCodeContent) FromJSONString(data string) error {
	return json.Unmarshal(common.StringToBytes(data), r)
}

type EmailRepo interface {
	SetCode(ctx context.Context, codeKey, content string, duration time.Duration) error
	GetCode(ctx context.Context, codeKey string) (content string, err error)
}
