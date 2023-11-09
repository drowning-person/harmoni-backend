package filesystem

import (
	"errors"
)

var (
	ErrUnknownPolicyType        = errors.New("unknown policy type")
	ErrFileSizeTooBig           = errors.New("file is too large")
	ErrFileExtensionNotAllowed  = errors.New("file type not allowed")
	ErrInsufficientCapacity     = errors.New("insufficient capacity")
	ErrIllegalObjectName        = errors.New("invalid object name")
	ErrClientCanceled           = errors.New("client canceled operation")
	ErrRootProtected            = errors.New("root protected")
	ErrInsertFileRecord         = errors.New("failed to create file record")
	ErrFileExisted              = errors.New("object existed")
	ErrFileUploadSessionExisted = errors.New("upload session existed")
	ErrPathNotExist             = errors.New("path not exist")
	ErrObjectNotExist           = errors.New("object not exist")
	ErrIO                       = errors.New("failed to read file data")
	ErrDBListObjects            = errors.New("failed to list object records")
	ErrDBDeleteObjects          = errors.New("failed to delete object records")
	ErrOneObjectOnly            = errors.New("you can only copy one object at the same time")
)
