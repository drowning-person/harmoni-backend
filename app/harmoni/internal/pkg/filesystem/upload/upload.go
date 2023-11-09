package upload

import (
	"encoding/json"
	"harmoni/app/harmoni/internal/pkg/common"
	"harmoni/app/harmoni/internal/pkg/filesystem/policy"
	"time"
)

// UploadCredential 返回给客户端的上传凭证
type UploadCredential struct {
	ChunkSize   uint64   `json:"chunkSize,omitempty"` // 分块大小，0 为部分快
	Expires     int64    `json:"expires"`             // 上传凭证过期时间， Unix 时间戳
	UploadURLs  []string `json:"uploadURLs,omitempty"`
	Credential  string   `json:"credential,omitempty"`
	UploadID    string   `json:"uploadID,omitempty"`
	Callback    string   `json:"callback,omitempty"` // 回调地址
	Path        string   `json:"path,omitempty"`     // 存储路径
	AccessKey   string   `json:"ak,omitempty"`
	KeyTime     string   `json:"keyTime,omitempty"` // COS用有效期
	Policy      string   `json:"policy,omitempty"`
	CompleteURL string   `json:"completeURL,omitempty"`
}

// UploadSession 上传会话
type UploadSession struct {
	UID            uint // 发起者
	FileID         int64
	SliceNum       uint64
	VirtualPath    string     // 用户文件路径，不含文件名
	Key            string     // 文件名
	Size           uint64     // 文件大小
	SavePath       string     // 物理存储路径，包含物理文件名
	LastModified   *time.Time // 可选的文件最后修改日期
	Policy         *policy.Policy
	Callback       string // 回调 URL 地址
	CallbackSecret string // 回调 URL
	UploadURL      string
	UploadID       string // 上传会话 GUID
	Credential     string
}

// UploadCallback 上传回调正文
type UploadCallback struct {
	PicInfo string `json:"pic_info"`
}

// GeneralUploadCallbackFailed 存储策略上传回调失败响应
type GeneralUploadCallbackFailed struct {
	Error string `json:"error"`
}

func (u *UploadSession) ToJSON() []byte {
	data, _ := json.Marshal(u)
	return data
}

func (u *UploadSession) FromJSON(data []byte) *UploadSession {
	json.Unmarshal(data, u)
	return u
}

func (u *UploadSession) FromJSONString(data string) *UploadSession {
	json.Unmarshal(common.StringToBytes(data), u)
	return u
}
