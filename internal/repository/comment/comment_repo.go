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

var _ commententity.CommentRepository = (*CommentRepo)(nil)

type CommentRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewCommentRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) *CommentRepo {
	return &CommentRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger.With("module", "repository/comment"),
	}
}

func (r *CommentRepo) Create(ctx context.Context, comment *commententity.Comment) error {
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

func (r *CommentRepo) GetByCommentID(ctx context.Context, commentID int64) (*commententity.Comment, bool, error) {
	comment := &commententity.Comment{}
	err := r.db.WithContext(ctx).Where("comment_id = ?", commentID).First(comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return comment, true, nil
}

func (r *CommentRepo) GetLikeCount(ctx context.Context, commentID int64) (int64, bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("comment").
		Select([]string{"like_count"}).
		Where("comment_id = ?", commentID).Scan(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, false, nil
		}
		return 0, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return count, true, nil
}

func (r *CommentRepo) UpdateLikeCount(ctx context.Context, commentID int64, count int64) error {
	if err := r.db.WithContext(ctx).
		Table("comment").
		Where("comment_id = ?", commentID).
		Update("like_count", count).Error; err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *CommentRepo) GetPage(ctx context.Context, commentQuery *commententity.CommentQuery) (paginator.Page[commententity.Comment], error) {
	commentPage := paginator.Page[commententity.Comment]{CurrentPage: commentQuery.Page, PageSize: commentQuery.PageSize}
	db := r.db.WithContext(ctx).Where("object_id = ? AND root_id = ?", commentQuery.ObjectID, commentQuery.RootID)

	if commentQuery.UserID != 0 {
		db = db.Where("user_id = ?", commentQuery.UserID)
	}

	switch commentQuery.QueryCond {
	case "newest":
		db.Order("created_at DESC")
	case "score":
		db.Order("like_count")
	default:
		db.Order("created_at DESC")
	}

	err := commentPage.SelectPages(db)
	if err != nil {
		return paginator.Page[commententity.Comment]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return commentPage, nil
}
