package cron

import (
	"context"
	"encoding/json"
	"harmoni/internal/conf"
	"harmoni/internal/entity"
	"harmoni/internal/entity/like"
	"harmoni/internal/pkg/queue/rabbitmq"
	"harmoni/internal/usecase"
	"time"

	"github.com/go-co-op/gocron"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type ScheduledTaskManager struct {
	producers   *rabbitmq.RabbitMqSender
	scheduler   *gocron.Scheduler
	likeUsecase *usecase.LikeUsecase
	logger      *zap.SugaredLogger
}

// NewScheduledTaskManager new scheduled task manager
func NewScheduledTaskManager(
	conf *conf.MessageQueue,
	likeUsecase *usecase.LikeUsecase,
	logger *zap.SugaredLogger,
) (*ScheduledTaskManager, func(), error) {
	s := gocron.NewScheduler(time.Local)
	manager := &ScheduledTaskManager{
		scheduler:   s,
		likeUsecase: likeUsecase,
		logger:      logger,
	}

	rabbitConf := rabbitmq.RabbitConf{
		Username: conf.RabbitMQ.Username,
		Password: conf.RabbitMQ.Password,
		Host:     conf.RabbitMQ.Host,
		Port:     conf.RabbitMQ.Port,
		VHost:    conf.RabbitMQ.VHost,
	}

	admin := rabbitmq.MustNewAdmin(rabbitConf)
	err := admin.DeclareExchange(rabbitmq.ExchangeConf{
		ExchangeName: entity.LikeExchange,
		Type:         amqp.ExchangeDirect,
		Durable:      true,
	}, nil)
	if err != nil {
		return nil, func() {}, err
	}
	err = admin.DeclareQueue(rabbitmq.QueueConf{
		Name:    entity.LikeQueue,
		Durable: true,
	}, nil)
	if err != nil {
		return nil, func() {}, err
	}
	err = admin.Bind(entity.LikeQueue, entity.LikeBindKey, entity.LikeExchange, false, nil)
	if err != nil {
		return nil, func() {}, err
	}

	manager.producers = rabbitmq.MustNewSender(rabbitmq.RabbitSenderConf{
		RabbitConf:  rabbitConf,
		ContentType: "application/json",
	})

	return manager, func() { manager.Stop() }, nil
}

func (s *ScheduledTaskManager) likeCountTask() {
	s.logger.Debug("start save like counts to DB")
	for _, likeType := range like.LikeTypeList {
		ctx := context.Background()
		counts, err := s.likeUsecase.BatchLikeCount(ctx, likeType)
		if err != nil {
			s.logger.Errorf("send like count msg to mq failed: %s", err)
			return
		}
		if len(counts) == 0 {
			continue
		}

		likeMsg := &like.LikeMessage{
			Type:     like.LikeCountMessage,
			LikeType: likeType,
			CountMessage: &like.CountMessage{
				Counts: counts,
			},
		}
		data, _ := json.Marshal(likeMsg)
		err = s.producers.Send(ctx, entity.LikeExchange, entity.LikeBindKey, data)
		if err != nil {
			s.logger.Errorf("send like count msg to mq failed: %s", err)
		}
	}
}

func (s *ScheduledTaskManager) Start() {
	s.scheduler.Every("5m").Do(s.likeCountTask)
	s.scheduler.StartAsync()
}

func (s *ScheduledTaskManager) Stop() {
	s.scheduler.Stop()
	err := s.producers.Close()
	if err != nil {
		s.logger.Errorf("stop cron failed: %s", err)
	}
}
