package comment

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/entity/unique"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"html"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type commentRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewCommentRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) commententity.CommentRepository {
	return &commentRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger,
	}
}

func (r *commentRepo) Create(ctx context.Context, comment *commententity.Comment) error {
	var err error
	comment.CommentID, err = r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return err
	}

	comment.Content = html.EscapeString(comment.Content)

	err = r.db.WithContext(ctx).Create(comment).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *commentRepo) GetPage(ctx context.Context, commentQuery *commententity.CommentQuery) (paginator.Page[commententity.Comment], error) {
	commentPage := paginator.Page[commententity.Comment]{CurrentPage: commentQuery.Page, PageSize: commentQuery.PageSize}
	db := r.db.WithContext(ctx).Where("object_id = ? AND root_id = ?", commentQuery.ObjectID, commentQuery.RootID)

	if commentQuery.UserID != 0 {
		db = db.Where("user_id = ?", commentQuery.UserID)
	}

	switch commentQuery.QueryCond {
	case "newest":
		db.Order("created_at DESC")
	default:
		db.Order("created_at DESC")
	}

	err := commentPage.SelectPages(db)
	if err != nil {
		return paginator.Page[commententity.Comment]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return commentPage, nil
}
