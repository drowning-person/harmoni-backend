package subscriber

import (
	"harmoni/app/harmoni/internal/infrastructure/config"

	"github.com/ThreeDotsLabs/watermill/message"
	"go.uber.org/zap"
)

func NewSubscriber(conf *config.MessageQueue, suffix string, logger *zap.Logger) (message.Subscriber, error) {
	var (
		sub message.Subscriber
		err error
	)
	switch {
	case conf.RabbitMQ != nil:
		sub, err = NewAMQPSubscriber(conf.RabbitMQ, suffix, logger)
		if err != nil {
			return nil, err
		}
	}
	return sub, nil
}
