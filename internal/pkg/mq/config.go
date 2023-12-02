package mq

import (
	"fmt"
	"harmoni/internal/conf"
)

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

func (c *MessageQueue) FromConfig(conf *conf.MessageQueue) *MessageQueue {
	c.RabbitMQ = &RabbitMQConf{
		Username: conf.RabbitMq.GetUsername(),
		Password: conf.RabbitMq.GetPassword(),
		Host:     conf.RabbitMq.GetHost(),
		Port:     int(conf.RabbitMq.GetPort()),
		VHost:    conf.RabbitMq.GetVhost(),
	}
	return c
}
