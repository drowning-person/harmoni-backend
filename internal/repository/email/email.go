package email

import (
	"context"
	emailrepo "harmoni/internal/entity/email"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"time"

	"github.com/redis/go-redis/v9"
)

// emailRepo email repository
type emailRepo struct {
	rdb *redis.Client
}

// NewEmailRepo new repository
func NewEmailRepo(rdb *redis.Client) emailrepo.EmailRepo {
	return &emailRepo{
		rdb: rdb,
	}
}

// SetCode The email code is used to verify that the link in the message is out of date
func (e *emailRepo) SetCode(ctx context.Context, codeKey, content string, duration time.Duration) error {
	err := e.rdb.Set(ctx, codeKey, content, duration).Err()
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// GetCode get the code
func (e *emailRepo) GetCode(ctx context.Context, codeKey string) (content string, err error) {
	content, err = e.rdb.Get(ctx, codeKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		err = errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
