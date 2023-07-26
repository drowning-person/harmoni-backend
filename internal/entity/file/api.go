package file

type AvatarUploadRequest struct {
	UserID   int64  `json:"-"`
	FileName string `json:"file_name,omitempty"`
	Size     int64  `json:"size,omitempty"`
	Content  []byte `json:"content,omitempty"`
}

type AvatarUploadResponse struct {
	Link string `json:"link,omitempty"`
}
