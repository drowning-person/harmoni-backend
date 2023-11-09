package email

import (
	"context"
	"encoding/json"
	"harmoni/app/harmoni/internal/pkg/common"
	"time"
)

type EmailCodeContent struct {
	Code        string `json:"code,omitempty"`
	LastReqTime int64  `json:"last_req_time,omitempty"`
}

type EmailType uint8

const (
	BindEmail EmailType = iota + 1
	ChangeEmail
	ChangePassword
	ResetPassword
)

func (r *EmailCodeContent) ToJSONString() string {
	codeBytes, _ := json.Marshal(r)
	return common.BytesToString(codeBytes)
}

func (r *EmailCodeContent) FromJSONString(data string) error {
	return json.Unmarshal(common.StringToBytes(data), r)
}

type EmailRepo interface {
	SetCode(ctx context.Context, codeKey, content string, duration time.Duration) (bool, error)
	GetCode(ctx context.Context, codeKey string) (content string, err error)
	DelCode(ctx context.Context, codeKey string) error
}
