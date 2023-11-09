package email

import (
	"context"
	emailentity "harmoni/app/harmoni/internal/entity/email"
	"harmoni/app/harmoni/internal/pkg/errorx"
	"harmoni/app/harmoni/internal/pkg/reason"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ emailentity.EmailRepo = (*EmailRepo)(nil)

// EmailRepo email repository
type EmailRepo struct {
	rdb *redis.Client
}

// NewEmailRepo new repository
func NewEmailRepo(rdb *redis.Client) *EmailRepo {
	return &EmailRepo{
		rdb: rdb,
	}
}

func (e *EmailRepo) DelCode(ctx context.Context, codeKey string) error {
	err := e.rdb.Del(ctx, codeKey).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

// SetCode The email code is used to verify that the link in the message is out of date
func (e *EmailRepo) SetCode(ctx context.Context, codeKey, content string, duration time.Duration) (bool, error) {
	err := e.rdb.SetArgs(ctx, codeKey, content, redis.SetArgs{
		Mode: "NX",
		TTL:  duration,
	}).Err()
	if err != nil {
		if err == redis.Nil {
			return true, nil
		}
		return false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return false, nil
}

// GetCode get the code
func (e *EmailRepo) GetCode(ctx context.Context, codeKey string) (content string, err error) {
	content, err = e.rdb.Get(ctx, codeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		err = errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
