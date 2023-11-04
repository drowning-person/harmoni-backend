package iface

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, value interface{}) error
	Close() error
}
