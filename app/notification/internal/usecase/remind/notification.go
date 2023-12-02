package remind

import (
	"context"
	"harmoni/app/notification/internal/entity/notifyconfig"
	"harmoni/app/notification/internal/entity/remind"
	"harmoni/internal/types/iface"

	"github.com/go-kratos/kratos/v2/log"
)

type RemindUsecase struct {
	nr     remind.RemindRepository
	cr     notifyconfig.NotifyConfigRepository
	tx     iface.Transaction
	logger *log.Helper
}

func NewRemindUsecase(
	cr notifyconfig.NotifyConfigRepository,
	nr remind.RemindRepository,
	tx iface.Transaction,
	logger log.Logger,
) *RemindUsecase {
	return &RemindUsecase{
		cr:     cr,
		nr:     nr,
		tx:     tx,
		logger: log.NewHelper(log.With(logger, "module", "usecase/notification")),
	}
}

func (u *RemindUsecase) Create(ctx context.Context, req *remind.CreateReq) error {
	config, err := u.cr.Get(ctx, req.Action, req.ObjectType)
	if err != nil {
		return err
	}
	remind := remind.Remind{}
	remind.BuildContent(config)
	req.Content = remind.Content
	return u.tx.ExecTx(ctx, func(ctx context.Context) error {
		return u.nr.Create(ctx, req)
	})
}

func (u *RemindUsecase) List(ctx context.Context, req *remind.ListReq) ([]*remind.Remind, error) {
	return u.nr.List(ctx, req)
}

func (u *RemindUsecase) Count(ctx context.Context, req *remind.CountReq) (int64, error) {
	return u.nr.Count(ctx, req)
}
