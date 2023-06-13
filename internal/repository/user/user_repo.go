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

var _ userentity.UserRepository = (*UserRepo)(nil)

func userVerifyKey(userID int64, verifyType userentity.VerifyType, actionType accountentity.AccountActionType) string {
	return fmt.Sprintf("%s%d:%s%d:%d", userKeyPrefix, userID, verifyKey, verifyType, actionType)
}

type UserRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewUserRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) *UserRepo {
	return &UserRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger.With("module", "repository/user"),
	}
}

func (r *UserRepo) Create(ctx context.Context, user *userentity.User) (err error) {
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

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*userentity.User, bool, error) {
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

func (r *UserRepo) GetByUserID(ctx context.Context, userID int64) (*userentity.User, bool, error) {
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

func (r *UserRepo) GetByUserIDs(ctx context.Context, userIDs []int64) ([]userentity.User, error) {
	users := make([]userentity.User, 0, 8)
	err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&users).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return users, nil
}

// GetPage get user page TODO: Add Condition
func (r *UserRepo) GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[userentity.User], error) {
	userPage := paginator.Page[userentity.User]{CurrentPage: int64(pageNum), PageSize: int64(pageSize)}
	err := userPage.SelectPages(r.db.WithContext(ctx))
	if err != nil {
		return paginator.Page[userentity.User]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return userPage, nil
}

func (r *UserRepo) GetModifyStaus(ctx context.Context, userID int64, verifyType userentity.VerifyType, actionType accountentity.AccountActionType) (userentity.ModifyStatus, error) {
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

func (r *UserRepo) SetModifyStatus(ctx context.Context, userID int64, status userentity.ModifyStatus, verifyType userentity.VerifyType, actionType accountentity.AccountActionType, statusKeepTime time.Duration) error {
	key := userVerifyKey(userID, verifyType, actionType)
	if _, err := r.rdb.Set(ctx, key, uint8(status), statusKeepTime).Result(); err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *UserRepo) ModifyPassword(ctx context.Context, user *userentity.User) error {
	err := r.db.WithContext(ctx).Model(user).Where("id = ?", user.ID).Update("password", user.Password).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *UserRepo) ModifyEmail(ctx context.Context, user *userentity.User) error {
	err := r.db.WithContext(ctx).Model(user).Where("id = ?", user.ID).Update("email", user.Email).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *UserRepo) GetLikeCount(ctx context.Context, userID int64) (int64, bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("user").
		Select([]string{"like_count"}).
		Where("user_id = ?", userID).Scan(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, false, nil
		}
		return 0, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return count, true, nil
}

func (r *UserRepo) UpdateLikeCount(ctx context.Context, userID int64, count int64) error {
	if err := r.db.WithContext(ctx).
		Table("user").
		Where("user_id = ?", userID).
		Update("like_count", count).Error; err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
