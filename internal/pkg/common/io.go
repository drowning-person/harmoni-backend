package common

import (
	"io"
	"mime/multipart"
	"os"
)

func ConvertMultipartFile(mf *multipart.FileHeader) (io.ReadSeekCloser, string, int64, error) {
	f, err := mf.Open()
	if err != nil {
		return nil, "", 0, err
	}
	return f, mf.Filename, mf.Size, nil
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
