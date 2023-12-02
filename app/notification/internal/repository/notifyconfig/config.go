package notifyconfig

import (
	"context"
	"errors"
	"harmoni/app/notification/internal/entity/notifyconfig"
	"harmoni/app/notification/internal/infrastructure/po/notification"
	"harmoni/internal/pkg/data"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/types/action"
	"harmoni/internal/types/object"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

var _ notifyconfig.NotifyConfigRepository = (*NotifyConfigRepo)(nil)

type NotifyConfigRepo struct {
	db     *data.DB
	logger *log.Helper
}

func NewNotifyConfigRepo(db *data.DB, logger log.Logger) *NotifyConfigRepo {
	return &NotifyConfigRepo{
		db:     db,
		logger: log.NewHelper(log.With(logger, "module", "usecase/notifyconfig")),
	}
}

func (r *NotifyConfigRepo) Get(ctx context.Context, action action.Action, objectType object.ObjectType) (*notifyconfig.NotifyConfig, error) {
	config := notification.NotifyConfig{}
	err := r.db.DB(ctx).
		Where("action = ? AND object_type = ?", action, objectType).
		First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return config.ToDomain(), nil
}
