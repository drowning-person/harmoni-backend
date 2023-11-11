package post

import (
	"context"
	"harmoni/app/harmoni/internal/entity/paginator"
	postentity "harmoni/app/harmoni/internal/entity/post"
	postreltagentity "harmoni/app/harmoni/internal/entity/post_rel_tag"
	tagentity "harmoni/app/harmoni/internal/entity/tag"
	"harmoni/app/harmoni/internal/entity/unique"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/errorx"
	"html"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ postentity.PostRepository = (*PostRepo)(nil)

type PostRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	tagRepo      tagentity.TagRepository
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewPostRepo(
	db *gorm.DB,
	rdb *redis.Client,
	tagRepo tagentity.TagRepository,
	uniqueIDRepo unique.UniqueIDRepo,
	logger *zap.SugaredLogger) *PostRepo {
	return &PostRepo{
		db:           db,
		rdb:          rdb,
		tagRepo:      tagRepo,
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

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&post).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		// add relationship between post and tags
		err = r.associateTags(ctx, tx, post.PostID, post.TagIDs)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepo) Delete(ctx context.Context, postID int64) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Where("post_id = ?", postID).
			Delete(&postentity.Post{}).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		err = r.removeAllTagsFromPost(ctx, tx, postID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
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

func (r *PostRepo) GetByUserID(ctx context.Context, userID int64, queryCond *postentity.PostQuery) (paginator.Page[postentity.Post], error) {
	postPage := paginator.Page[postentity.Post]{CurrentPage: queryCond.Page, PageSize: queryCond.PageSize}
	db := r.db.WithContext(ctx).Where("author_id = ?", userID).Order("created_at DESC")
	err := postPage.SelectPages(db)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return paginator.Page[postentity.Post]{}, errorx.NotFound(reason.PostNotFound)
		}
		return paginator.Page[postentity.Post]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return postPage, nil
}

func (r *PostRepo) GetByUserIDs(ctx context.Context, userIDs []int64, queryCond *postentity.PostQuery) (paginator.Page[postentity.Post], error) {
	postPage := paginator.Page[postentity.Post]{CurrentPage: queryCond.Page, PageSize: queryCond.PageSize}
	db := r.db.WithContext(ctx).Where("author_id IN ?", userIDs).Order("created_at DESC")
	err := postPage.SelectPages(db)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return paginator.Page[postentity.Post]{}, errorx.NotFound(reason.PostNotFound)
		}
		return paginator.Page[postentity.Post]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return postPage, nil
}

func (r *PostRepo) BatchByIDs(ctx context.Context, postIDs []int64) ([]postentity.Post, error) {
	posts := make([]postentity.Post, 0, len(postIDs))
	if err := r.db.WithContext(ctx).
		Select([]string{"author_id", "post_id", "title", "content", "status"}).
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

func (r *PostRepo) GetPage(ctx context.Context, queryCond *postentity.PostQuery) (paginator.Page[postentity.Post], error) {
	db := r.db.WithContext(ctx)
	if len(queryCond.AuthorIDs) != 0 {
		if len(queryCond.AuthorIDs) == 1 {
			db = db.Where("author_id in ?", queryCond.AuthorIDs)
		} else {
			db = db.Where("author_id = ?", queryCond.AuthorIDs[0])
		}
		db = db.Order("created_at DESC")
	} else {
		if queryCond.TagID != 0 {
			db = db.Joins("INNER JOIN post_tags AS pt ON pt.post_id = post.post_id").Where("pt.tag_id = ?", queryCond.TagID)
		}
		switch queryCond.QueryCond {
		case postentity.PostOrderByCreatedTime:
			db = db.Order("created_at DESC")
		case postentity.PostOrderByLike:
			db = db.Order("like_count DESC")
		default:
			db = db.Order("created_at DESC")
		}
	}

	postPage := paginator.Page[postentity.Post]{CurrentPage: queryCond.Page, PageSize: queryCond.PageSize}
	err := postPage.SelectPages(db)
	if err != nil {
		return paginator.Page[postentity.Post]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return postPage, nil
}

func (r *PostRepo) associateTags(ctx context.Context, tx *gorm.DB, postID int64, tagIDs []int64) error {
	if len(tagIDs) == 0 {
		return nil
	}

	prt := make([]postreltagentity.PostTag, len(tagIDs))
	prtID, err := r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return err
	}

	for i, tagID := range tagIDs {
		prt[i] = postreltagentity.PostTag{
			PostTagID: prtID,
			PostID:    postID,
			TagID:     tagID,
		}
	}

	err = tx.Create(&prt).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *PostRepo) GetPostsByTagID(ctx context.Context, tagID int64) ([]postentity.Post, error) {
	posts := []postentity.Post{}
	err := r.db.Table(postreltagentity.TableName).
		Where("tag_id = ?", tagID).
		Find(posts).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return posts, nil
}

/* func (r *PostRepo) RemoveTagsFromPost(ctx context.Context, postID int64, tagIDs []int64) error {
	err := r.db.Table(postreltagentity.TableName).
		Where("post_id = ? and tag_id in ?", postID, tagIDs).
		Delete(&postreltagentity.PostTag{}).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
} */

func (r *PostRepo) removeAllTagsFromPost(ctx context.Context, tx *gorm.DB, postID int64) error {
	err := r.db.Table(postreltagentity.TableName).
		Where("post_id = ?", postID).
		Delete(&postreltagentity.PostTag{}).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}
