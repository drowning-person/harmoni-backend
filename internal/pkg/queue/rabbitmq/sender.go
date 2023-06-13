package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Sender interface {
		Send(exchange string, routeKey string, msg []byte) error
	}

	RabbitMqSender struct {
		conn        *amqp.Connection
		channel     *amqp.Channel
		ContentType string
	}
)

func MustNewSender(rabbitMqConf RabbitSenderConf) *RabbitMqSender {
	sender := &RabbitMqSender{ContentType: rabbitMqConf.ContentType}
	conn, err := amqp.Dial(getRabbitURL(rabbitMqConf.RabbitConf))
	if err != nil {
		panic(fmt.Errorf("failed to connect rabbitmq, error: %v", err))
	}

	sender.conn = conn
	channel, err := sender.conn.Channel()
	if err != nil {
		panic(fmt.Errorf("failed to open a channel, error: %v", err))
	}

	sender.channel = channel
	return sender
}

func (q *RabbitMqSender) Send(ctx context.Context, exchange string, routeKey string, msg []byte) error {
	return q.channel.PublishWithContext(
		ctx,
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  q.ContentType,
			DeliveryMode: 2,
			Body:         msg,
		},
	)
}

func (q *RabbitMqSender) Close() error {
	return q.conn.Close()
}
