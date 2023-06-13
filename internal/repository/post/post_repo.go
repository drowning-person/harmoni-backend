package post

import (
	"context"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
	"harmoni/internal/entity/unique"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"html"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ postentity.PostRepository = (*PostRepo)(nil)

type PostRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewPostRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) *PostRepo {
	return &PostRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger.With("module", "repository/post"),
	}
}

func (r *PostRepo) Create(ctx context.Context, post *postentity.Post) (err error) {
	post.PostID, err = r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return err
	}

	post.Title = html.EscapeString(post.Title)
	post.Content = html.EscapeString(post.Content)

	err = r.db.WithContext(ctx).Create(post).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *PostRepo) GetBasicInfoByPostID(ctx context.Context, postID int64) (*postentity.Post, bool, error) {
	post := &postentity.Post{}
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).First(post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return post, true, nil
}

func (r *PostRepo) GetByPostID(ctx context.Context, postID int64) (*postentity.Post, bool, error) {
	post, exist, err := r.GetBasicInfoByPostID(ctx, postID)
	if err != nil {
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return post, exist, nil
}

func (r *PostRepo) BatchByIDs(ctx context.Context, postIDs []int64) ([]postentity.Post, error) {
	posts := make([]postentity.Post, 0, len(postIDs))
	if err := r.db.WithContext(ctx).Where("post_id IN ?", postIDs).Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.NotFound(reason.PostNotFound)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return posts, nil
}

func (r *PostRepo) BatchBasicInfoByIDs(ctx context.Context, postIDs []int64) ([]postentity.Post, error) {
	posts := make([]postentity.Post, 0, len(postIDs))
	if err := r.db.WithContext(ctx).
		Select([]string{"author_id", "post_id", "title", "content"}).
		Where("post_id IN ?", postIDs).Find(&posts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorx.NotFound(reason.PostNotFound)
		}
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return posts, nil
}

func (r *PostRepo) GetLikeCount(ctx context.Context, postID int64) (int64, bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Table("post").
		Select([]string{"like_count"}).
		Where("post_id = ?", postID).Scan(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, false, nil
		}
		return 0, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return count, true, nil
}

func (r *PostRepo) UpdateLikeCount(ctx context.Context, postID int64, count int64) error {
	if err := r.db.WithContext(ctx).
		Table("post").
		Where("post_id = ?", postID).
		Update("like_count", count).Error; err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *PostRepo) GetPage(ctx context.Context, pageSize, pageNum int64, orderCond string) (paginator.Page[postentity.Post], error) {
	db := r.db.WithContext(ctx)

	switch orderCond {
	case "newest":
		db.Order("created_at DESC")
	case "score":
		db.Order("like_count")
	default:
		db.Order("created_at DESC")
	}

	postPage := paginator.Page[postentity.Post]{CurrentPage: pageNum, PageSize: pageSize}
	err := postPage.SelectPages(db)
	if err != nil {
		return paginator.Page[postentity.Post]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return postPage, nil
}
