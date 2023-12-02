package remind

import (
	"context"
	v1 "harmoni/app/notification/api/http/v1/notification"
	"harmoni/app/notification/internal/entity/remind"
	usecaseremind "harmoni/app/notification/internal/usecase/remind"

	"github.com/go-kratos/kratos/v2/log"
)

type RemindService struct {
	ru     *usecaseremind.RemindUsecase
	logger *log.Helper
}

func NewRemindService(
	ru *usecaseremind.RemindUsecase,
	logger log.Logger,
) *RemindService {
	return &RemindService{
		ru:     ru,
		logger: log.NewHelper(log.With(logger, "module", "service/remind")),
	}
}

func (s *RemindService) UnreadCount(ctx context.Context, req *v1.UnReadRequest) (*v1.UnReadResponse, error) {
	count, err := s.ru.Count(ctx, &remind.CountReq{
		UserID: req.UserID,
		Action: req.Action,
		UnRead: true,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UnReadResponse{Count: count}, nil
}
