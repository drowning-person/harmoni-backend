package user

import (
	"context"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/entity/unique"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewUserRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) userentity.UserRepository {
	return &userRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger,
	}
}

func (r *userRepo) Create(ctx context.Context, user *userentity.User) (err error) {
	user.UserID, err = r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*userentity.User, bool, error) {
	user := &userentity.User{}
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return user, true, nil
}

func (r *userRepo) GetByUserID(ctx context.Context, userID int64) (*userentity.User, bool, error) {
	user := &userentity.User{}
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return user, true, nil
}

// GetPage get user page TODO: Add Condition
func (r *userRepo) GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[userentity.User], error) {
	userPage := paginator.Page[userentity.User]{CurrentPage: int64(pageNum), PageSize: int64(pageSize)}
	err := userPage.SelectPages(r.db.WithContext(ctx))
	if err != nil {
		return paginator.Page[userentity.User]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return userPage, nil
}
