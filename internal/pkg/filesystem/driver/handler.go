package driver

import (
	"context"
	"errors"
	"harmoni/internal/pkg/filesystem/fsctx"
	"harmoni/internal/pkg/filesystem/response"
	"harmoni/internal/pkg/filesystem/upload"
)

var (
	ErrInvalidParts   = errors.New("one or more of the specified parts could not be found")
	ErrEntityTooSmall = errors.New("your proposed upload is smaller than the minimum allowed object size")
)

type Part struct {
	PartNumber   int
	Size         int
	LastModified string
	ETag         string
}

type ByPartNumber []Part

func (p ByPartNumber) Len() int {
	return len(p)
}

func (p ByPartNumber) Less(i, j int) bool {
	return p[i].PartNumber < p[j].PartNumber
}

func (p ByPartNumber) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Handler interface {
	// 上传文件, dst为文件存储路径，size 为文件大小。上下文关闭
	// 时，应取消上传并清理临时文件
	Put(ctx context.Context, file fsctx.FileHeader) (string, error)

	// 合并分块
	Merge(ctx context.Context, file fsctx.FileHeader, parts []Part) (string, error)

	// 列出分块
	ListParts(ctx context.Context, file fsctx.FileHeader, maxParts int64, offset int64) ([]Part, error)

	// 删除分块
	DeleteParts(ctx context.Context, file fsctx.FileHeader) error

	// 删除一个或多个给定路径的文件，返回删除失败的文件路径列表及错误
	Delete(ctx context.Context, files []string) ([]string, error)

	// 获取文件内容
	Get(ctx context.Context, path string) (response.RSCloser, error)

	// 获取外链/下载地址，
	// url - 站点本身地址,
	// isDownload - 是否直接下载
	Source(ctx context.Context, rootURL string, path string, ttl int64, isDownload bool, speed int) (string, error)

	// Token 获取有效期为ttl的上传凭证和签名
	Token(ctx context.Context, ttl int64, uploadSession *upload.UploadSession, file fsctx.FileHeader) (*upload.UploadCredential, error)

	// CancelToken 取消已经创建的有状态上传凭证
	CancelToken(ctx context.Context, uploadSession *upload.UploadSession) error
}
