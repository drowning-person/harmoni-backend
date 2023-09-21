package response

import (
	"io"
)

// RSCloser 存储策略适配器返回的文件流，有些策略需要带有Closer
type RSCloser interface {
	io.ReadSeeker
	io.Closer
}
