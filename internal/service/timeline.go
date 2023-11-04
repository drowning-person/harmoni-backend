package service

import (
	"context"
	postentity "harmoni/internal/entity/post"
	timelineentity "harmoni/internal/entity/timeline"
	"harmoni/internal/usecase/timeline"

	"go.uber.org/zap"
)

type TimeLineService struct {
	tc     *timeline.TimeLinePullUsecase
	logger *zap.SugaredLogger
}

func NewTimeLineService(
	tc *timeline.TimeLinePullUsecase,
	logger *zap.SugaredLogger,
) *TimeLineService {
	return &TimeLineService{
		tc:     tc,
		logger: logger,
	}
}

func (s *TimeLineService) GetUserTimeLine(ctx context.Context, req *timelineentity.GetUserTimeLineRequest) (*timelineentity.GetUserTimeLineReply, error) {
	timeline, err := s.tc.GetTimeLineByUserID(ctx, req.AuthorID, req.UserID, &postentity.PostQuery{PageCond: req.PageCond})
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &timelineentity.GetUserTimeLineReply{
		Page: *timeline,
	}, nil
}

func (s *TimeLineService) GetHomeTimeLine(ctx context.Context, req *timelineentity.GetHomeTimeLineRequest) (*timelineentity.GetHomeTimeLineReply, error) {
	timeline, err := s.tc.GetTimeLine(ctx, req.UserID, &postentity.PostQuery{PageCond: req.PageCond})
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &timelineentity.GetHomeTimeLineReply{
		Page: *timeline,
	}, nil
}
