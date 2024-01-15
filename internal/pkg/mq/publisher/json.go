package publisher

import (
	"context"
	"encoding/json"
	"fmt"

	"harmoni/internal/pkg/mq"
	"harmoni/internal/pkg/watermillkratos"
	"harmoni/internal/types/iface"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"google.golang.org/protobuf/proto"
)

var ProviderSetPublisher = wire.NewSet(
	NewJSONPublisher,
	wire.Bind(new(iface.Publisher), new(*JSONPublisher)),
)

func NewPublisher(conf *mq.MessageQueue, logger log.Logger) (message.Publisher, error) {
	var (
		pub message.Publisher
		err error
	)
	switch {
	case conf.RabbitMQ != nil:
		amqpConfig := amqp.NewDurablePubSubConfig(conf.RabbitMQ.BuildURL(), nil)
		pub, err = amqp.NewPublisher(amqpConfig, watermillkratos.NewLogger(
			log.NewHelper(log.With(logger, "module", "publisher/like", "service", "like")), "msg",
		))
		if err != nil {
			return nil, err
		}
	}

	return pub, nil
}

var _ iface.Publisher = (*JSONPublisher)(nil)

type JSONPublisher struct {
	message.Publisher
}

func NewJSONPublisher(publisher message.Publisher) *JSONPublisher {
	return &JSONPublisher{
		Publisher: publisher,
	}
}

func (p *JSONPublisher) Publish(ctx context.Context, topic string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), data)
	return p.Publisher.Publish(topic, msg)
}

type ProtoBufPublisher struct {
	message.Publisher
}

func NewProtoBufPublisher(publisher message.Publisher) *ProtoBufPublisher {
	return &ProtoBufPublisher{
		Publisher: publisher,
	}
}

func (p *ProtoBufPublisher) Publish(ctx context.Context, topic string, value interface{}) error {
	v, ok := value.(proto.Message)
	if !ok {
		return fmt.Errorf("value is not a proto.Message: %T", value)
	}
	data, err := proto.Marshal(v)
	if err != nil {
		return err
	}
	msg := message.NewMessage(watermill.NewUUID(), data)
	return p.Publisher.Publish(topic, msg)
}
