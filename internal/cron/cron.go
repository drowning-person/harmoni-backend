package cron

import (
	"context"
	"harmoni/internal/entity/like"
	"harmoni/internal/infrastructure/config"
	eventlike "harmoni/internal/types/events/like"
	"harmoni/internal/types/iface"
	likeusecase "harmoni/internal/usecase/like"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

var _ iface.Executor = (*ScheduledTaskManager)(nil)

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

	return manager, func() { manager.Shutdown() }, nil
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

func (s *ScheduledTaskManager) Start() error {
	_, err := s.scheduler.Every(s.conf.DatabaseSyncInterval).Do(s.likeCountTask)
	if err != nil {
		return err
	}
	s.scheduler.StartAsync()
	return nil
}

func (s *ScheduledTaskManager) Shutdown() error {
	s.scheduler.Stop()
	err := s.publisher.Close()
	if err != nil {
		s.logger.Errorf("stop cron failed: %s", err)
	}
	return err
}
