package filesystem

import (
	"context"
	"errors"
	"harmoni/internal/pkg/filesystem/driver"
	"harmoni/internal/pkg/filesystem/driver/local"
	"harmoni/internal/pkg/filesystem/policy"
	"harmoni/internal/pkg/filesystem/response"
	"sync"

	"github.com/redis/go-redis/v9"
)

type FileSystem struct {
	// 操作文件使用的存储策略
	Policy *policy.Policy
	// 互斥锁
	Lock sync.Mutex

	/*
	   钩子函数
	*/
	Hooks map[string][]Hook

	rdb redis.UniversalClient

	/*
	   文件系统处理适配器
	*/
	Handler driver.Handler
}

// NewFileSystem 初始化一个文件系统
func NewFileSystem(policy *policy.Policy, rdb redis.UniversalClient) (*FileSystem, error) {
	fs := &FileSystem{}
	fs.Policy = policy
	fs.rdb = rdb
	// 分配存储策略适配器
	err := fs.DispatchHandler()

	return fs, err
}

// DispatchHandler 根据存储策略分配文件适配器
func (fs *FileSystem) DispatchHandler() error {
	if fs.Policy == nil {
		return errors.New("未设置存储策略")
	}
	policyType := fs.Policy.Type
	currentPolicy := fs.Policy

	switch policyType {
	case "mock", "anonymous":
		return nil
	case "local":
		fs.Handler = local.NewDriver(currentPolicy.BucketName, currentPolicy, fs.rdb)
		return nil
	default:
		return ErrUnknownPolicyType
	}
}

func (fs *FileSystem) GetContentByPath(ctx context.Context, path string) (response.RSCloser, error) {
	// 获取文件流
	rs, err := fs.Handler.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
