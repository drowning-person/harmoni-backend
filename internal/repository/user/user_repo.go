package user

import (
	"context"
	"fmt"
	accountentity "harmoni/internal/entity/account"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/entity/unique"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	userKeyPrefix = "user:"
	verifyKey     = "verify:"
)

func userVerifyKey(userID int64, verifyType userentity.VerifyType, actionType accountentity.AccountActionType) string {
	return fmt.Sprintf("%s%d:%s%d:%d", userKeyPrefix, userID, verifyKey, verifyType, actionType)
}

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

func (r *userRepo) GetModifyStaus(ctx context.Context, userID int64, verifyType userentity.VerifyType, actionType accountentity.AccountActionType) (userentity.ModifyStatus, error) {
	key := userVerifyKey(userID, verifyType, actionType)
	statusTmp, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return userentity.NotVerifiedEmailOrPhone, nil
		}
		return userentity.NotVerifiedEmailOrPhone, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	status, err := strconv.ParseUint(statusTmp, 10, 8)
	if err != nil {
		r.logger.Errorf("convert redis string to ModifyStatus failed: %s", err)
	}

	return userentity.ModifyStatus(status), nil
}

func (r *userRepo) SetModifyStatus(ctx context.Context, userID int64, status userentity.ModifyStatus, verifyType userentity.VerifyType, actionType accountentity.AccountActionType, statusKeepTime time.Duration) error {
	key := userVerifyKey(userID, verifyType, actionType)
	if _, err := r.rdb.Set(ctx, key, uint8(status), statusKeepTime).Result(); err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userRepo) ModifyPassword(ctx context.Context, user *userentity.User) error {
	err := r.db.WithContext(ctx).Model(user).Where("id = ?", user.ID).Update("password", user.Password).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *userRepo) ModifyEmail(ctx context.Context, user *userentity.User) error {
	err := r.db.WithContext(ctx).Model(user).Where("id = ?", user.ID).Update("email", user.Email).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
