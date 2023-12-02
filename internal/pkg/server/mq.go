package server

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kratos/kratos/v2/transport"
)

var _ transport.Server = (*MQServer)(nil)

type MQServer struct {
	*message.Router
}

func (r *MQServer) Start(ctx context.Context) error {
	return r.Run(ctx)
}

func (r *MQServer) Stop(context.Context) error {
	return r.Close()
}

func NewMQServer(r *message.Router) *MQServer {
	return &MQServer{
		Router: r,
	}
}
