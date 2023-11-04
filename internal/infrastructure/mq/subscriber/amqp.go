package subscriber

import (
	"harmoni/internal/infrastructure/config"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
)

func NewAMQPSubscriber(conf *config.RabbitMQConf, suffix string) (*amqp.Subscriber, error) {
	amqpConfig := amqp.NewDurablePubSubConfig(conf.BuildURL(), amqp.GenerateQueueNameTopicNameWithSuffix(suffix))
	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}
