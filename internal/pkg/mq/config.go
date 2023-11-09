package mq

import "fmt"

type MessageQueue struct {
	RabbitMQ *RabbitMQConf `mapstructure:"rabbitmq,omitempty"`
}

type RabbitMQConf struct {
	Username string `mapstructure:"username,omitempty"`
	Password string `mapstructure:"password,omitempty"`
	Host     string `mapstructure:"host,omitempty"`
	Port     int    `mapstructure:"port,omitempty"`
	VHost    string `mapstructure:"vhost,omitempty"`
}

func (c *RabbitMQConf) BuildURL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", c.Username, c.Password, c.Host, c.Port, c.VHost)
}
