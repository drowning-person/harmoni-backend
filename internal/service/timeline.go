package service

import (
	"context"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
	timelineentity "harmoni/internal/entity/timeline"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type TimeLineService struct {
	tc     *usecase.TimeLinePullUsecase
	logger *zap.SugaredLogger
}

func NewTimeLineService(
	tc *usecase.TimeLinePullUsecase,
	logger *zap.SugaredLogger,
) *TimeLineService {
	return &TimeLineService{
		tc:     tc,
		logger: logger,
	}
}

func (s *TimeLineService) GetUserTimeLine(ctx context.Context, req *timelineentity.GetUserTimeLineRequest) (*timelineentity.GetUserTimeLineReply, error) {
	timeline, err := s.tc.GetTimeLineByUserID(ctx, req.UserID, &postentity.PostQuery{PageCond: req.PageCond})
	if err != nil {
		return nil, err
	}

	res := paginator.Page[postentity.PostDetail]{
		CurrentPage: timeline.CurrentPage,
		PageSize:    timeline.PageSize,
		Total:       timeline.Total,
		Pages:       timeline.Pages,
		Data:        make([]postentity.PostDetail, 0, len(timeline.Data)),
	}

	for _, post := range timeline.Data {
		res.Data = append(res.Data, postentity.ConvertPostToDisplayDetail(&post))
	}

	return &timelineentity.GetUserTimeLineReply{
		Page: res,
	}, nil
}

func (s *TimeLineService) GetHomeTimeLine(ctx context.Context, req *timelineentity.GetHomeTimeLineRequest) (*timelineentity.GetHomeTimeLineReply, error) {
	timeline, err := s.tc.GetTimeLine(ctx, req.UserID, &postentity.PostQuery{PageCond: req.PageCond})
	if err != nil {
		return nil, err
	}

	res := paginator.Page[postentity.PostDetail]{
		CurrentPage: timeline.CurrentPage,
		PageSize:    timeline.PageSize,
		Total:       timeline.Total,
		Pages:       timeline.Pages,
		Data:        make([]postentity.PostDetail, 0, len(timeline.Data)),
	}

	for _, post := range timeline.Data {
		res.Data = append(res.Data, postentity.ConvertPostToDisplayDetail(&post))
	}

	return &timelineentity.GetHomeTimeLineReply{
		Page: res,
	}, nil
}
