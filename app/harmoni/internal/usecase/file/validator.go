package file

import (
	"context"
	"harmoni/app/harmoni/internal/pkg/common"
	"strings"
)

// 文件/路径名保留字符
var reservedCharacter = []string{"\\", "?", "*", "<", "\"", ":", ">", "/", "|"}

// ValidateLegalName 验证文件名/文件夹名是否合法
func (u *FileUseCase) ValidateLegalName(ctx context.Context, name string) bool {
	// 是否包含保留字符
	for _, value := range reservedCharacter {
		if strings.Contains(name, value) {
			return false
		}
	}

	// 是否超出长度限制
	if len(name) >= 256 {
		return false
	}

	// 是否为空限制
	if len(name) == 0 {
		return false
	}

	// 结尾不能是空格
	if strings.HasSuffix(name, " ") {
		return false
	}

	return true
}

// ValidateFileSize 验证上传的文件大小是否超出限制
func (u *FileUseCase) ValidateFileSize(ctx context.Context, size uint64) bool {
	if u.fs.Policy.MaxSize == 0 {
		return true
	}
	return size <= u.fs.Policy.MaxSize
}

// ValidateExtension 验证文件扩展名
func (u *FileUseCase) ValidateExtension(ctx context.Context, fileName string) bool {
	// 不需要验证
	if len(u.fs.Policy.OptionsSerialized.FileType) == 0 {
		return true
	}

	return common.IsInExtensionList(u.fs.Policy.OptionsSerialized.FileType, fileName)
}
