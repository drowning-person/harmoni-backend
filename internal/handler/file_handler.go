package handler

import (
	fileentity "harmoni/internal/entity/file"
	"harmoni/internal/pkg/common"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx/fiberx"
	"harmoni/internal/pkg/middleware"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/service"
	"path/filepath"

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
	req.Content, _, req.Size, err = common.ConvertMultipartFile(header)
	if err != nil {
		h.logger.Errorln(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.UploadAvatar(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *FileHandler) GetFileContent(c *fiber.Ctx) error {
	req := fileentity.GetFileContentRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	reply, err := h.ffs.GetFileContent(c.Context(), &req)
	if err != nil {
		return fiberx.HandleResponse(c, err, nil)
	}

	return c.Type(filepath.Ext(req.FilePath)).SendStream(reply.Content)
}

func (h *FileHandler) UploadPrepare(c *fiber.Ctx) error {
	req := fileentity.UploadPrepareRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.UploadPrepare(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *FileHandler) UploadPart(c *fiber.Ctx) error {
	req := fileentity.UploadMultipartRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}

	header, err := c.FormFile("file")
	if err != nil {
		h.logger.Errorln(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError), nil)
	}
	req.Content, _, req.Size, err = common.ConvertMultipartFile(header)
	if err != nil {
		h.logger.Errorln(err)
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.UploadPart(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *FileHandler) UploadComplete(c *fiber.Ctx) error {
	req := fileentity.UploadMultipartCompleteRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.UploadComplete(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *FileHandler) ListParts(c *fiber.Ctx) error {
	req := fileentity.ListPartsRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.ListParts(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}

func (h *FileHandler) AbortMultipartUpload(c *fiber.Ctx) error {
	req := fileentity.AbortMultipartUploadRequest{}
	if err := fiberx.ParseAndCheck(c, &req); err != nil {
		return fiberx.HandleResponse(c, errorx.BadRequest(reason.RequestFormatError).WithMsg(err.Error()), nil)
	}
	req.UserID = middleware.GetClaimsFromCtx(c.UserContext()).UserID

	reply, err := h.ffs.AbortMultipartUpload(c.Context(), &req)
	return fiberx.HandleResponse(c, err, reply)
}
