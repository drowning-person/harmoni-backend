package handler

import (
	"bytes"
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"
	"io"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type FileHandler struct {
	ffs    *service.FileService
	logger *zap.SugaredLogger
}

func NewFileHandler(
	ffs *service.FileService,
	logger *zap.SugaredLogger,
) *FileHandler {
	return &FileHandler{
		ffs:    ffs,
		logger: logger.With("module", "handler/file"),
	}
}

func (h *FileHandler) UploadAvatar(c *fiber.Ctx) error {
	req := fileentity.AvatarUploadRequest{}
	header, err := c.FormFile("file")
	if err != nil {
		h.logger.Errorln(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError), nil)
	}

	req.FileName = header.Filename
	req.Size = header.Size
	file, err := header.Open()
	if err != nil {
		h.logger.Errorln(err)
		return fiberx.HandleResponse(c, err, nil)
	}

	content := &bytes.Buffer{}
	written, err := io.Copy(content, file)
	if written != header.Size {
		h.logger.Errorf("want write %v bytes, but %v bytes", header.Size, written)
		return fiberx.HandleResponse(c, err, nil)
	}
	req.Content = content.Bytes()
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.UploadAvatar(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}
