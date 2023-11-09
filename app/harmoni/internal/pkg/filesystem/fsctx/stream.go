package fsctx

import (
	"errors"
	"harmoni/app/harmoni/internal/pkg/filesystem/policy"
	"io"
	"time"
)

type WriteMode int

const (
	Overwrite WriteMode = 0x00001
	// Append 只适用于本地策略
	Append WriteMode = 0x00002
	Nop    WriteMode = 0x00004
)

type UploadTaskInfo struct {
	Size            uint64
	IsPart          bool
	PartNumber      int
	Hash            string
	Ext             string
	FileName        string
	Mode            WriteMode
	Metadata        map[string]string
	LastModified    *time.Time
	SavePath        string
	UploadSessionID *string
	Model           interface{}
	ETag            string
	Policy          *policy.Policy
}

// FileHeader 上传来的文件数据处理器
type FileHeader interface {
	io.Reader
	io.Closer
	io.Seeker
	Info() *UploadTaskInfo
	SetSize(uint64)
	SetModel(fileModel interface{})
	Seekable() bool
}

// FileStream 用户传来的文件
type FileStream struct {
	ETag         string
	IsPart       bool
	PartNumber   int
	Mode         WriteMode
	Hash         string
	Ext          string
	LastModified *time.Time
	Metadata     map[string]string
	File         io.ReadCloser
	Seeker       io.Seeker
	Size         uint64
	Name         string
	SavePath     string
	UploadID     string
	Model        interface{}
	Src          string
	Policy       *policy.Policy
}

func (file *FileStream) Read(p []byte) (n int, err error) {
	if file.File != nil {
		return file.File.Read(p)
	}

	return 0, io.EOF
}

func (file *FileStream) Close() error {
	if file.File != nil {
		return file.File.Close()
	}

	return nil
}

func (file *FileStream) Seek(offset int64, whence int) (int64, error) {
	if file.Seekable() {
		return file.Seeker.Seek(offset, whence)
	}

	return 0, errors.New("no seeker")
}

func (file *FileStream) Seekable() bool {
	return file.Seeker != nil
}

func (file *FileStream) Info() *UploadTaskInfo {
	return &UploadTaskInfo{
		Size:            file.Size,
		ETag:            file.ETag,
		IsPart:          file.IsPart,
		PartNumber:      file.PartNumber,
		Hash:            file.Hash,
		Ext:             file.Ext,
		FileName:        file.Name,
		Mode:            file.Mode,
		Metadata:        file.Metadata,
		LastModified:    file.LastModified,
		SavePath:        file.SavePath,
		UploadSessionID: &file.UploadID,
		Model:           file.Model,
		Policy:          file.Policy,
	}
}

func (file *FileStream) SetSize(size uint64) {
	file.Size = size
}

func (file *FileStream) SetModel(fileModel interface{}) {
	file.Model = fileModel
}
