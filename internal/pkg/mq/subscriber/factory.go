package subscriber

import (
	"harmoni/internal/pkg/mq"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewSubscriber(conf *mq.MessageQueue, suffix string, logger watermill.LoggerAdapter) (message.Subscriber, error) {
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
