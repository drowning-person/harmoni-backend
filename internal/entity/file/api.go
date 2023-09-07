package file

import (
	"harmoni/internal/pkg/filesystem/driver"
	"harmoni/internal/pkg/filesystem/response"
	"harmoni/internal/pkg/filesystem/upload"
	"io"
)

type AvatarUploadRequest struct {
	UserID   int64             `json:"-"`
	FileName string            `json:"file_name,omitempty"`
	Size     int64             `json:"size,omitempty"`
	Content  io.ReadSeekCloser `json:"-"`
}

type AvatarUploadResponse struct {
	Link string `json:"link,omitempty"`
}

type CreateUploadSessionRequest struct {
	UserID int64  `json:"-"`
	Size   uint64 `json:"size" validate:"required"`
	Name   string `json:"name" validate:"required"`
}

type CreateUploadSessionResponse struct {
	Credential *upload.UploadCredential
}

type GetFileContentRequest struct {
	FilePath string `params:"filepath"`
}

type GetFileContentResponse struct {
	Content response.RSCloser
}

type UploadObjectRequest struct {
	UserID   int64             `json:"-"`
	FileName string            `json:"file_name,omitempty"`
	Size     int64             `json:"size,omitempty"`
	Content  io.ReadSeekCloser `json:"-"`
}

type UploadObjectResponse struct {
	Location string `json:"location"`
}

type IsObjectUploadedRequest struct {
	Hash string `query:"hash" validate:"required"`
}

type IsObjectUploadedResponse struct {
	Location string `json:"location"`
}

type UploadPrepareRequest struct {
	MD5    string `json:"md5" validate:"required"`
	Key    string `json:"key" validate:"required"`
	UserID int64  `json:"-"`
}

type UploadPrepareResponse struct {
	Credential *upload.UploadCredential
}

type UploadMultipartRequest struct {
	UserID     int64             `json:"-"`
	Key        string            `json:"key" form:"key" validate:"required"`
	UploadID   string            `json:"upload_id" validate:"required" form:"upload_id"`
	PartNumber int               `json:"part_number" validate:"required" form:"part_number"`
	Size       int64             `json:"-"`
	Content    io.ReadSeekCloser `json:"-"`
}

type UploadMultipartResponse struct {
	Etag string `json:"etag"` // MD5
}

type FilePart struct {
	PartNumber int    `json:"part_number"`
	Etag       string `json:"etag"`
}

type Parts []FilePart

func (f Parts) ToDriver() []driver.Part {
	ps := make([]driver.Part, len(f))
	for i, part := range f {
		ps[i] = driver.Part{
			PartNumber: part.PartNumber,
			ETag:       part.Etag,
		}
	}
	return ps
}

type UploadMultipartCompleteRequest struct {
	UserID    int64      `json:"-"`
	Key       string     `json:"key" validate:"required"`
	UploadID  string     `json:"upload_id" validate:"required"`
	FileParts []FilePart `json:"file_parts" validate:"required"`
}

type UploadMultipartCompleteResponse struct {
	Etag     string `json:"etag"`
	Location string `json:"location"` // 文件地址
}

type ListPartsRequest struct {
	UserID           int64  `json:"-,omitempty"`
	Key              string `params:"key,omitempty" validate:"required"`
	UploadID         string `query:"upload_id,omitempty" validate:"required"`
	MaxParts         int64  `query:"max_parts,omitempty"`
	PartNumberMarker int64  `query:"part_number_marker,omitempty"`
}

type ListPartsResponse struct {
	Parts []driver.Part `json:"parts,omitempty"`
}

type AbortMultipartUploadRequest struct {
	UserID   int64  `json:"-"`
	Key      string `json:"key" validate:"required"`
	UploadID string `json:"upload_id" validate:"required"`
}

type AbortMultipartUploadResponse struct {
}
