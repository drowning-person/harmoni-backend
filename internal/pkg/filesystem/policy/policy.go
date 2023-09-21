package policy

type FileType uint8

const (
	Image FileType = iota + 1
	Video
	Audio
)

const (
	// 图片拓展
	ExtJPEG = ".jpeg"
	ExtJPG  = ".jpg"
	ExtPNG  = ".png"
	ExtGIF  = ".gif"
	ExtBMP  = ".bmp"
	ExtTIFF = ".tiff"
	ExtWebP = ".webp"
	ExtSVG  = ".svg"

	// 音频拓展
	ExtMP3  = ".mp3"
	ExtWAV  = ".wav"
	ExtOGG  = ".ogg"
	ExtFLAC = ".flac"
	ExtAAC  = ".aac"

	// 视频拓展
	ExtMP4  = ".mp4"
	ExtAVI  = ".avi"
	ExtMKV  = ".mkv"
	ExtWMV  = ".wmv"
	ExtFLV  = ".flv"
	ExtMOV  = ".mov"
	ExtWebM = ".webm"
)

// Policy 存储策略
type Policy struct {
	Type       string
	BucketName string
	MaxSize    uint64
	// key type, value dir
	DirRule map[FileType]string
	// 数据库忽略字段
	OptionsSerialized PolicyOption
}

// PolicyOption 非公有的存储策略属性
type PolicyOption struct {
	// 允许的文件扩展名
	FileType  []string `json:"file_type"`
	ChunkSize uint64   `json:"chunk_size,omitempty"`
}

func (policy *Policy) GeneratePathByDirRule(ext string) string {
	switch ext {
	case ExtJPEG,
		ExtJPG,
		ExtPNG,
		ExtGIF,
		ExtBMP,
		ExtTIFF,
		ExtWebP,
		ExtSVG:
		return policy.DirRule[Image]
	case ExtMP3,
		ExtWAV,
		ExtOGG,
		ExtFLAC,
		ExtAAC:
		return policy.DirRule[Audio]
	case ExtMP4,
		ExtAVI,
		ExtMKV,
		ExtWMV,
		ExtFLV,
		ExtMOV,
		ExtWebM:
		return policy.DirRule[Video]
	}
	return ""
}
