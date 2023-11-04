package subscriber

import (
	"harmoni/internal/infrastructure/config"

	"github.com/ThreeDotsLabs/watermill/message"
)

func NewSubscriber(conf *config.MessageQueue, suffix string) (message.Subscriber, error) {
	var (
		sub message.Subscriber
		err error
	)
	switch {
	case conf.RabbitMQ != nil:
		sub, err = NewAMQPSubscriber(conf.RabbitMQ, suffix)
		if err != nil {
			return nil, err
		}
	}
	return sub, nil
}
