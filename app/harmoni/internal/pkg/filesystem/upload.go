package filesystem

import (
	"bytes"
	"context"
	"fmt"
	"harmoni/app/harmoni/internal/pkg/filesystem/fsctx"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Upload 上传文件
func (fs *FileSystem) Upload(ctx context.Context, file *fsctx.FileStream) (etag string, err error) {
	// 上传前的钩子
	err = fs.Trigger(ctx, "BeforeUpload", file)
	if err != nil {
		// request.BlackHole(file)
		return "", err
	}

	// 生成文件名和路径,
	var savePath string
	if file.SavePath == "" {
		if file.IsPart {
			savePath = fs.GenerateSavePathByChunk(ctx, file)
		} else {
			savePath = fs.GenerateSavePathByDirRule(ctx, file)
		}
		file.SavePath = savePath
	}

	// 保存文件
	if file.Mode&fsctx.Nop != fsctx.Nop {
		// 处理客户端未完成上传时，关闭连接
		go fs.CancelUpload(ctx, savePath, file)

		etag, err = fs.Handler.Put(ctx, file)
		if err != nil {
			fs.Trigger(ctx, "AfterUploadFailed", file)
			return "", err
		}
	}

	// 上传完成后的钩子
	err = fs.Trigger(ctx, "AfterUpload", file)

	if err != nil {
		// 上传完成后续处理失败
		followUpErr := fs.Trigger(ctx, "AfterValidateFailed", file)
		// 失败后再失败...
		if followUpErr != nil {
			return "", followUpErr
		}

		return "", err
	}

	return
}

func (fs *FileSystem) GenerateSavePathByDirRule(ctx context.Context, file fsctx.FileHeader) string {
	fileInfo := file.Info()
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fileInfo.Hash)
	buf.WriteString(fileInfo.Ext)
	return path.Join(
		fs.Policy.GeneratePathByDirRule(
			fileInfo.Ext,
		),
		buf.String(),
	)
}

func (fs *FileSystem) GenerateSavePathByChunk(ctx context.Context, file fsctx.FileHeader) string {
	fileInfo := file.Info()
	buf := bytes.NewBuffer(nil)
	buf.WriteString(strings.ReplaceAll(*fileInfo.UploadSessionID, "-", ""))
	buf.WriteString(fileInfo.Ext)
	if fileInfo.PartNumber != 0 {
		buf.WriteString(fmt.Sprintf("-%d", fileInfo.PartNumber))
	}
	return path.Join(
		fs.Policy.GeneratePathByDirRule(
			fileInfo.Ext,
		),
		buf.String(),
	)
}

// CancelUpload 监测客户端取消上传
func (fs *FileSystem) CancelUpload(ctx context.Context, path string, file fsctx.FileHeader) {
	var reqContext context.Context
	if fiberCtx, ok := ctx.Value(fsctx.FiberCtx).(*fiber.Ctx); ok {
		reqContext = fiberCtx.Context()
	} else if reqCtx, ok := ctx.Value(fsctx.HTTPCtx).(context.Context); ok {
		reqContext = reqCtx
	} else {
		return
	}

	<-reqContext.Done()
	select {
	case <-ctx.Done():
		// 客户端正常关闭，不执行操作
	default:
		// 客户端取消上传，删除临时文件
		if fs.Hooks["AfterUploadCanceled"] == nil {
			return
		}
		err := fs.Trigger(ctx, "AfterUploadCanceled", file)
		if err != nil {
			// util.Log().Debug("AfterUploadCanceled hook execution failed: %s", err)
		}
	}
}
