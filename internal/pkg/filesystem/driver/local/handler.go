package local

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"harmoni/internal/pkg/common"
	"harmoni/internal/pkg/filesystem/driver"
	"harmoni/internal/pkg/filesystem/fsctx"
	"harmoni/internal/pkg/filesystem/policy"
	"harmoni/internal/pkg/filesystem/response"
	"harmoni/internal/pkg/filesystem/upload"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	Perm   = 0744
	prefix = "local:"
)

func getStorePrefix(uploadID string) string {
	return prefix + uploadID
}

// LocalStorage 本地策略适配器
type LocalStorage struct {
	Policy *policy.Policy
	path   string
	rdb    redis.UniversalClient
}

type PartInfo struct {
	PartNumber   int
	Size         int
	LastModified string
	Etag         string
	Path         string
}

func (p *PartInfo) ToPart() driver.Part {
	return driver.Part{
		PartNumber:   p.PartNumber,
		Size:         p.Size,
		LastModified: p.LastModified,
		ETag:         p.Etag,
	}
}

func (p *PartInfo) ToJSON() string {
	data, _ := json.Marshal(p)
	return common.BytesToString(data)
}

func (p *PartInfo) FromJSONString(data string) *PartInfo {
	json.Unmarshal(common.StringToBytes(data), p)
	return p
}

type PartInfos []PartInfo

func (ps PartInfos) ToParts() []driver.Part {
	parts := make([]driver.Part, len(ps))
	for i, partInfo := range ps {
		parts[i] = partInfo.ToPart()
	}
	return parts
}

var _ driver.Handler = (*LocalStorage)(nil)

func NewDriver(path string, policy *policy.Policy, rdb redis.UniversalClient) *LocalStorage {
	return &LocalStorage{
		Policy: policy,
		path:   path,
		rdb:    rdb,
	}
}

func (ls LocalStorage) getTargetPath(path string) (string, error) {
	tmpPath := ls.path
	var err error
	if !filepath.IsAbs(tmpPath) {
		tmpPath, err = filepath.Abs(tmpPath)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(tmpPath, path), nil
}

// Get 获取文件内容
func (ls LocalStorage) Get(ctx context.Context, path string) (response.RSCloser, error) {
	// 打开文件
	dst, err := ls.getTargetPath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(dst)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Put 将文件流保存到指定目录
func (ls LocalStorage) Put(ctx context.Context, file fsctx.FileHeader) (string, error) {
	defer file.Close()
	fileInfo := file.Info()
	dst, err := ls.getTargetPath(fileInfo.SavePath)
	if err != nil {
		return "", err
	}

	// 如果非 Overwrite，则检查是否有重名冲突
	if fileInfo.Mode&fsctx.Overwrite != fsctx.Overwrite {
		if common.Exists(dst) {
			return "", errors.New("file with the same name existed or unavailable")
		}
	}

	// 如果目标目录不存在，创建
	basePath := filepath.Dir(dst)
	if !common.Exists(basePath) {
		err := os.MkdirAll(basePath, Perm)
		if err != nil {
			return "", err
		}
	}

	var (
		out *os.File
	)

	openMode := os.O_CREATE | os.O_RDWR
	if fileInfo.Mode&fsctx.Append == fsctx.Append {
		openMode |= os.O_APPEND
	} else {
		openMode |= os.O_TRUNC
	}

	out, err = os.OpenFile(dst, openMode, Perm)
	if err != nil {
		return "", err
	}
	defer out.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	n, err := out.Write(data)
	if err != nil {
		return "", err
	} else if n != len(data) {
		return "", fmt.Errorf("write data failed.want %d but %d", len(data), n)
	}

	hasher := md5.New()
	hasher.Write(data)
	hashed := hex.EncodeToString(hasher.Sum(nil))
	if fileInfo.IsPart {
		ls.rdb.ZAdd(ctx, getStorePrefix(*fileInfo.UploadSessionID), redis.Z{
			Score: float64(fileInfo.PartNumber),
			Member: (&PartInfo{
				PartNumber:   fileInfo.PartNumber,
				Size:         len(data),
				Path:         dst,
				LastModified: time.Now().Format("2006-01-02 15:04:05.000"),
				Etag:         hashed,
			}).ToJSON(),
		})
	}

	return hashed, nil
}

func (ls LocalStorage) Merge(ctx context.Context, file fsctx.FileHeader, parts []driver.Part) (string, error) {
	defer file.Close()
	fileInfo := file.Info()
	dst, err := ls.getTargetPath(fileInfo.SavePath)
	if err != nil {
		return "", err
	}

	openMode := os.O_CREATE | os.O_RDWR
	out, err := os.OpenFile(dst, openMode, Perm)
	if err != nil {
		return "", err
	}

	sort.Sort(driver.ByPartNumber(parts))
	partNumbers := make([]int, len(parts))
	for i, part := range parts {
		if part.PartNumber != i+1 {
			return "", errors.New("the parts are not successive")
		}
		partNumbers[i] = part.PartNumber
	}
	infos, err := ls.rdb.ZRangeByScore(ctx, getStorePrefix(*fileInfo.UploadSessionID), &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return "", err
	} else if len(infos) != len(parts) {
		return "", driver.ErrInvalidParts
	}

	hash := md5.New()
	var offset int64
	for _, partInfoStr := range infos {
		if partInfoStr == "" {
			return "", driver.ErrInvalidParts
		}
		partInfo := (&PartInfo{}).FromJSONString(partInfoStr)
		if partInfo.Etag != parts[partInfo.PartNumber-1].ETag {
			return "", driver.ErrInvalidParts
		} else if partInfo.Size < int(fileInfo.Policy.OptionsSerialized.ChunkSize) {
			return "", fmt.Errorf("%w. each part must be at least %d MB in size, except the last part", driver.ErrEntityTooSmall,
				fileInfo.Policy.OptionsSerialized.ChunkSize/1024/1024)
		}

		file, err := os.Open(partInfo.Path)
		if err != nil {
			return "", err
		}
		info, err := file.Stat()
		if err != nil {
			return "", err
		}
		data := make([]byte, info.Size())
		n, err := file.Read(data)
		if err != nil {
			return "", err
		} else if int64(n) != (info.Size()) {
			return "", fmt.Errorf("want %d bytes but read %d bytes", info.Size(), n)
		}
		hash.Write(data)
		n, err = out.WriteAt(data, offset)
		if err != nil {
			return "", err
		} else if int64(n) != (info.Size()) {
			return "", fmt.Errorf("want %d bytes but read %d bytes", info.Size(), n)
		}

		err = os.Remove(partInfo.Path)
		if err != nil {
			return "", err
		}
		offset += info.Size()
	}

	err = ls.rdb.Del(ctx, getStorePrefix(*fileInfo.UploadSessionID)).Err()
	if err != nil {
		return "", err
	}
	file.SetSize(uint64(offset))
	mergedHash := hash.Sum(nil)
	mergedHashString := hex.EncodeToString(mergedHash)
	return mergedHashString, nil
}

func (d LocalStorage) ListParts(ctx context.Context, file fsctx.FileHeader, maxParts int64, offset int64) ([]driver.Part, error) {
	fileInfo := file.Info()
	opt := &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: int64(offset),
	}
	if maxParts > 0 {
		opt.Count = int64(maxParts)
	}
	infos, err := d.rdb.ZRangeByScore(ctx, getStorePrefix(*fileInfo.UploadSessionID), opt).Result()
	if err != nil {
		return nil, err
	}

	partInfos := make([]PartInfo, len(infos))
	for i, partInfoStr := range infos {
		partInfos[i] = *(&PartInfo{}).FromJSONString(partInfoStr)
	}

	return PartInfos(partInfos).ToParts(), nil
}

func (ls LocalStorage) DeleteParts(ctx context.Context, file fsctx.FileHeader) error {
	fileInfo := file.Info()
	key := getStorePrefix(*fileInfo.UploadSessionID)
	infos, err := ls.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return err
	}

	for _, partInfoStr := range infos {
		if partInfoStr == "" {
			return driver.ErrInvalidParts
		}
		partInfo := (&PartInfo{}).FromJSONString(partInfoStr)
		err := os.Remove(partInfo.Path)
		if err != nil {
			return err
		}
	}
	err = ls.rdb.Unlink(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (ls LocalStorage) Truncate(ctx context.Context, src string, size uint64) error {
	// common.Log().Warning("Truncate file %q to [%d].", src, size)
	out, err := os.OpenFile(src, os.O_WRONLY, Perm)
	if err != nil {
		// common.Log().Warning("Failed to open file: %s", err)
		return err
	}

	defer out.Close()
	return out.Truncate(int64(size))
}

// Delete 删除一个或多个文件，
// 返回未删除的文件，及遇到的最后一个错误
func (ls LocalStorage) Delete(ctx context.Context, files []string) ([]string, error) {
	deleteFailed := make([]string, 0, len(files))
	var retErr error

	for _, value := range files {
		filePath := filepath.Join(ls.path, value)
		if common.Exists(filePath) {
			err := os.Remove(filePath)
			if err != nil {
				// common.Log().Warning("Failed to delete file: %s", err)
				retErr = err
				deleteFailed = append(deleteFailed, value)
			}
		}

		// 尝试删除文件的缩略图（如果有）
		// _ = os.Remove(common.RelativePath(value + model.GetSettingByNameWithDefault("thumb_file_suffix", "._thumb")))
	}

	return deleteFailed, retErr
}

// Source 获取外链URL
func (ls LocalStorage) Source(ctx context.Context, rootURL, path string, ttl int64, isDownload bool, speed int) (string, error) {
	var (
		signedURI *url.URL
		err       error
	)

	joined, err := url.JoinPath(rootURL, fmt.Sprintf("/api/v1/file/get/%s", path))
	if err != nil {
		return "", err
	}
	signedURI, err = url.Parse(joined)

	if err != nil {
		return "", fmt.Errorf("failed to sign url err: %s", err)
	}

	finalURL := signedURI.String()

	return finalURL, nil
}

// Token 获取上传策略和认证Token，本地策略直接返回空值
func (ls LocalStorage) Token(ctx context.Context, ttl int64, uploadSession *upload.UploadSession, file fsctx.FileHeader) (*upload.UploadCredential, error) {
	fileInfo := file.Info()
	dst, err := ls.getTargetPath(fileInfo.SavePath)
	if err != nil {
		return nil, err
	}
	fmt.Println(dst)
	if common.Exists(dst) {
		return nil, errors.New("placeholder file already exist")
	}

	return &upload.UploadCredential{
		UploadID: uploadSession.UploadID,
	}, nil
}

// 取消上传凭证
func (ls LocalStorage) CancelToken(ctx context.Context, uploadSession *upload.UploadSession) error {
	return nil
}
