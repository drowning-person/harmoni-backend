package mq

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Group struct {
	*message.Router
	Name string
	Sub  message.Subscriber
	Pub  message.Publisher
}

func (g *Group) Handle(topic string, f message.NoPublishHandlerFunc) *message.Handler {
	return g.Router.AddNoPublisherHandler(
		fmt.Sprintf("%s.%s", g.Name, topic),
		topic,
		g.Sub,
		f,
	)
}

func (g *Group) HandleAndPublish(topic, pubTopic string, f message.HandlerFunc) *message.Handler {
	return g.Router.AddHandler(
		fmt.Sprintf("%s.%s", g.Name, topic),
		topic,
		g.Sub,
		pubTopic,
		g.Pub,
		f,
	)
}
