package etcdx

import (
	"harmoni/internal/conf"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func NewETCDClient(
	conf *conf.ETCD,
	logger *zap.Logger,
) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   conf.GetAddr(),
		Username:    conf.GetUsername(),
		Password:    conf.GetPassword(),
		DialTimeout: 5 * time.Second,
	})
}
