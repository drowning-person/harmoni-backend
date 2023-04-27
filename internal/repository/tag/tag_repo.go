package tag

import (
	"context"
	"harmoni/internal/entity/paginator"
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/entity/unique"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"html"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type tagRepo struct {
	db           *gorm.DB
	rdb          *redis.Client
	uniqueIDRepo unique.UniqueIDRepo
	logger       *zap.SugaredLogger
}

func NewTagRepo(db *gorm.DB, rdb *redis.Client, uniqueIDRepo unique.UniqueIDRepo, logger *zap.SugaredLogger) tagentity.TagRepository {
	return &tagRepo{
		db:           db,
		rdb:          rdb,
		uniqueIDRepo: uniqueIDRepo,
		logger:       logger,
	}
}

func (r *tagRepo) Create(ctx context.Context, tag *tagentity.Tag) (err error) {
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

func (r *tagRepo) GetByTagID(ctx context.Context, tagID int64) (*tagentity.Tag, bool, error) {
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

func (r *tagRepo) GetByTagName(ctx context.Context, tagName string) (*tagentity.Tag, bool, error) {
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
func (r *tagRepo) GetPage(ctx context.Context, pageSize, pageNum int64) (paginator.Page[tagentity.Tag], error) {
	tagPage := paginator.Page[tagentity.Tag]{CurrentPage: pageNum, PageSize: pageSize}
	err := tagPage.SelectPages(r.db.WithContext(ctx))
	if err != nil {
		return paginator.Page[tagentity.Tag]{}, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	return tagPage, nil
}
