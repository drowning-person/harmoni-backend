package subscriber

import (
	"harmoni/internal/pkg/mq"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
)

func NewAMQPSubscriber(conf *mq.RabbitMQConf, suffix string, logger watermill.LoggerAdapter) (*amqp.Subscriber, error) {
	amqpConfig := amqp.NewDurablePubSubConfig(conf.BuildURL(), amqp.GenerateQueueNameTopicNameWithSuffix(suffix))
	amqpConfig.Consume.NoRequeueOnNack = true 
	subscriber, err := amqp.NewSubscriber(
		amqpConfig,
		logger,
	)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}
