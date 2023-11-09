package subscriber

import (
	"harmoni/app/harmoni/internal/infrastructure/config"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/garsue/watermillzap"
	"go.uber.org/zap"
)

func NewAMQPSubscriber(conf *config.RabbitMQConf, suffix string, logger *zap.Logger) (*amqp.Subscriber, error) {
	amqpConfig := amqp.NewDurablePubSubConfig(conf.BuildURL(), amqp.GenerateQueueNameTopicNameWithSuffix(suffix))
	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		watermillzap.NewLogger(logger),
	)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}
