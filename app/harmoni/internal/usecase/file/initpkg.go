package file

import (
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/pkg/filesystem/policy"
)

func NewPolicy(conf *config.FileStorage) *policy.Policy {
	return &policy.Policy{
		Type:       conf.Policy.Type,
		BucketName: conf.Policy.BucketName,
		MaxSize:    conf.Policy.MaxSize,
		// key type, value dir
		DirRule: conf.Policy.DirRule,
		OptionsSerialized: policy.PolicyOption{
			FileType:  conf.Policy.Option.FileType,
			ChunkSize: conf.Policy.Option.ChunkSize,
		},
	}
}
