package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	ConsumeHandle func(message string) error

	ConsumeHandler interface {
		Consume(message []byte) error
	}

	RabbitListener struct {
		conn         *amqp.Connection
		channel      *amqp.Channel
		forever      chan bool
		handler      ConsumeHandler
		queues       RabbitListenerConf
		consumedTags chan uint64 // 存储已消费的消息的 delivery tag
	}
)

func MustNewListener(listenerConf RabbitListenerConf, handler ConsumeHandler) *RabbitListener {
	listener := &RabbitListener{
		queues:       listenerConf,
		handler:      handler,
		forever:      make(chan bool),
		consumedTags: make(chan uint64, 1),
	}

	conn, err := amqp.Dial(getRabbitURL(listenerConf.RabbitConf))
	if err != nil {
		log.Fatalf("failed to connect rabbitmq, error: %v", err)
	}

	listener.conn = conn
	channel, err := listener.conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}

	listener.channel = channel

	go listener.startAckWorker()

	return listener
}

func (q *RabbitListener) startAckWorker() {
	for {
		select {
		case tag := <-q.consumedTags:
			q.ack(tag)
		case <-q.forever:
			return
		}
	}
}

func (q *RabbitListener) ack(tag uint64) {
	q.channel.Ack(tag, false)
}

func (q *RabbitListener) Start() {
	for _, que := range q.queues.ListenerQueues {
		msg, err := q.channel.Consume(
			que.Name,
			"",
			que.AutoAck,
			que.Exclusive,
			que.NoLocal,
			que.NoWait,
			nil,
		)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		go func() {
			for d := range msg {
				if err := q.handler.Consume(d.Body); err != nil {
					log.Printf("Error on consuming: %s, error: %v", string(d.Body), err)
				}
				q.consumedTags <- d.DeliveryTag
			}
		}()
	}

	<-q.forever
}

func (q *RabbitListener) Stop() {
	q.channel.Close()
	q.conn.Close()
	close(q.forever)
}
