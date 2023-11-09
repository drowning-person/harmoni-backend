package tag

import (
	"context"
	"harmoni/app/harmoni/internal/entity/paginator"
	postreltagentity "harmoni/app/harmoni/internal/entity/post_rel_tag"
	tagentity "harmoni/app/harmoni/internal/entity/tag"
	"harmoni/app/harmoni/internal/entity/unique"
	"harmoni/app/harmoni/internal/pkg/errorx"
	"harmoni/app/harmoni/internal/pkg/reason"
	"html"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ tagentity.TagRepository = (*TagRepo)(nil)

type TagRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewTagRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) *TagRepo {
	return &TagRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger.With("module", "repository/tag"),
	}
}

func (r *TagRepo) Create(ctx context.Context, tag *tagentity.Tag) (err error) {
	tag.TagID, err = r.uniqueIDRepo.GenUniqueID(ctx)
	if err != nil {
		return err
	}

	r.logger.Debugf("Create %#v", tag)

	tag.TagName = html.EscapeString(tag.TagName)
	tag.Introduction = html.EscapeString(tag.Introduction)

	r.logger.Debugf("Create %#v", tag)

	err = r.db.WithContext(ctx).Create(tag).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return nil
}

func (r *TagRepo) GetByTagID(ctx context.Context, tagID int64) (*tagentity.Tag, bool, error) {
	tag := &tagentity.Tag{}
	err := r.db.WithContext(ctx).Where("tag_id = ?", tagID).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tag, true, nil
}

func (r *TagRepo) GetByTagIDs(ctx context.Context, tagIDs []int64) ([]tagentity.Tag, error) {
	tags := make([]tagentity.Tag, 0, 8)
	err := r.db.WithContext(ctx).Where("tag_id IN ?", tagIDs).Find(&tags).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tags, nil
}

func (r *TagRepo) GetByTagName(ctx context.Context, tagName string) (*tagentity.Tag, bool, error) {
	tag := &tagentity.Tag{}
	err := r.db.WithContext(ctx).Where("tag_name = ?", tagName).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil
		}
		return nil, false, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tag, true, nil
}

// GetPage get tag page TODO: Add Condition
func (r *TagRepo) GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[tagentity.Tag], error) {
	tagPage := paginator.Page[tagentity.Tag]{CurrentPage: pageNum, PageSize: pageSize}
	err := tagPage.SelectPages(r.db.WithContext(ctx))
	if err != nil {
		return paginator.Page[tagentity.Tag]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tagPage, nil
}

func (r *TagRepo) GetTagsByPostID(ctx context.Context, postID int64) ([]tagentity.Tag, error) {
	tags := []tagentity.Tag{}
	err := r.db.Table(postreltagentity.TableName).
		Select("tag.tag_id", "tag.tag_name").
		Where("post_id = ?", postID).
		Joins("JOIN tag on post_tags.tag_id = tag.tag_id").
		Find(&tags).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tags, nil
}
