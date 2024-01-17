package remind

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"

	v1 "harmoni/app/notification/api/http/v1/notification"
	"harmoni/app/notification/internal/entity/remind"
	usecaseremind "harmoni/app/notification/internal/usecase/remind"
	"harmoni/internal/pkg/httpx"
	"harmoni/internal/pkg/paginator"
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

func (s *RemindService) ListRemind(ctx context.Context, req *v1.ListRemindRequest) (*v1.ListRemindResponse, error) {
	list, err := s.ru.List(ctx, &remind.ListReq{
		UserID: req.UserID,
		Action: req.Action,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ListRemindResponse{
		PageResp: httpx.PageResp{
			Total: list.Total,
			Page: httpx.Page{
				PageNum:  list.CurrentPage,
				PageSize: list.PageSize,
			},
		},
		Reminds: lo.Map(list.Data, func(remind *remind.Remind, _ int) *v1.Remind {
			return ConverRemindToResp(remind)
		}),
	}, nil
}

func (s *RemindService) LikeDetail(ctx context.Context, req *v1.LikeRemindDetailRequest) (*v1.LikeRemindDetailResponse, error) {
	remindSenders, err := s.ru.ListRemindSenders(ctx, &remind.ListRemindSendersReq{
		PageRequest: &paginator.PageRequest{
			Num:  int64(req.PageNum),
			Size: int64(req.PageSize),
		},
		RemindID: req.RemindID,
		Action:   req.Action,
	})
	if err != nil {
		return nil, err
	}
	resp := &v1.LikeRemindDetailResponse{
		PageResp: httpx.PageResp{
			Total: remindSenders.Total,
			Page: httpx.Page{
				PageNum:  remindSenders.CurrentPage,
				PageSize: remindSenders.PageSize,
			},
		},
		Items: lo.Map(remindSenders.Data, func(r *remind.RemindSender, _ int) *v1.LikeRemindDetailItem {
			return &v1.LikeRemindDetailItem{
				User:      r.Sender,
				CreatedAt: r.CreatedAt,
			}
		}),
	}
	return resp, nil
}
