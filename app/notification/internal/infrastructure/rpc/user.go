package rpc

import (
	"context"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcdclient "go.etcd.io/etcd/client/v3"
)

func NewUsergRPC(
	client *etcdclient.Client,
) (v1.UserClient, error) {
	reg := etcd.New(client)
	endpoint := "discovery:///harmoni"
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(endpoint), grpc.WithDiscovery(reg))
	if err != nil {
		return nil, err
	}
	return v1.NewUserClient(conn), nil
}
