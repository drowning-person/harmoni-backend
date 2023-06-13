package rabbitmq

import "fmt"

type RabbitConf struct {
	Username string
	Password string
	Host     string
	Port     int
	VHost    string
}

type RabbitListenerConf struct {
	RabbitConf
	ListenerQueues []ConsumerConf
}

type ConsumerConf struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	// Set to true, which means that messages sent by producers in the same connection
	// cannot be delivered to consumers in this connection.
	NoLocal bool
	// Whether to block processing
	NoWait bool
}

type RabbitSenderConf struct {
	RabbitConf
	ContentType string // MIME content type
}

type QueueConf struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

type ExchangeConf struct {
	ExchangeName string
	Type         string // exchange type
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	Queues       []QueueConf
}

func getRabbitURL(rabbitConf RabbitConf) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitConf.Username, rabbitConf.Password,
		rabbitConf.Host, rabbitConf.Port, rabbitConf.VHost)
}
