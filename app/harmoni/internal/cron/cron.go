package cron

import (
	"context"
	"harmoni/app/harmoni/internal/entity/like"
	"harmoni/app/harmoni/internal/infrastructure/config"
	eventlike "harmoni/app/harmoni/internal/types/events/like"
	"harmoni/app/harmoni/internal/types/iface"
	likeusecase "harmoni/app/harmoni/internal/usecase/like"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-kratos/kratos/v2/transport"
	"go.uber.org/zap"
)

var _ transport.Server = (*ScheduledTaskManager)(nil)

type ScheduledTaskManager struct {
	conf        *config.Like
	publisher   iface.Publisher
	scheduler   *gocron.Scheduler
	likeUsecase *likeusecase.LikeUsecase
	logger      *zap.SugaredLogger
}

// NewScheduledTaskManager new scheduled task manager
func NewScheduledTaskManager(
	conf *config.Like,
	publisher iface.Publisher,
	likeUsecase *likeusecase.LikeUsecase,
	logger *zap.SugaredLogger,
) (*ScheduledTaskManager, func(), error) {
	s := gocron.NewScheduler(time.Local)
	manager := &ScheduledTaskManager{
		conf:        conf,
		scheduler:   s,
		publisher:   publisher,
		likeUsecase: likeUsecase,
		logger:      logger,
	}

	return manager, func() { manager.Stop(context.Background()) }, nil
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

		likeMsg := &eventlike.LikeStoreMessage{
			BaseMessage: eventlike.BaseMessage{
				LikeType: likeType.ToEventLikeType(),
			},
			Counts: counts,
		}
		err = s.publisher.Publish(ctx, eventlike.TopicLikeStore, likeMsg)
		if err != nil {
			s.logger.Errorf("send like count msg to mq failed: %s", err)
		}
	}
}

func (s *ScheduledTaskManager) Start(context.Context) error {
	_, err := s.scheduler.Every(s.conf.DatabaseSyncInterval).Do(s.likeCountTask)
	if err != nil {
		return err
	}
	s.scheduler.StartAsync()
	return nil
}

func (s *ScheduledTaskManager) Stop(context.Context) error {
	s.scheduler.Stop()
	err := s.publisher.Close()
	if err != nil {
		s.logger.Errorf("stop cron failed: %s", err)
	}
	return err
}
