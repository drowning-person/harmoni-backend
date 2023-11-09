package tag

import (
	"context"
	"harmoni/app/harmoni/internal/entity/paginator"
	tagentity "harmoni/app/harmoni/internal/entity/tag"

	"go.uber.org/zap"
)

type TagUseCase struct {
	tagRepo tagentity.TagRepository
	logger  *zap.SugaredLogger
}

func NewTagUseCase(tagRepo tagentity.TagRepository, logger *zap.SugaredLogger) *TagUseCase {
	return &TagUseCase{
		tagRepo: tagRepo,
		logger:  logger,
	}
}

func (u *TagUseCase) Create(ctx context.Context, tag *tagentity.Tag) (*tagentity.Tag, error) {
	err := u.tagRepo.Create(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, err
}

func (u *TagUseCase) GetByTagID(ctx context.Context, tagID int64) (*tagentity.Tag, bool, error) {
	return u.tagRepo.GetByTagID(ctx, tagID)
}

func (u *TagUseCase) GetByTagIDs(ctx context.Context, tagIDs []int64) ([]tagentity.Tag, error) {
	return u.tagRepo.GetByTagIDs(ctx, tagIDs)
}

func (u *TagUseCase) GetByTagName(ctx context.Context, name string) (*tagentity.Tag, bool, error) {
	return u.tagRepo.GetByTagName(ctx, name)
}

func (u *TagUseCase) GetPage(ctx context.Context, pageSize int64, pageNum int64) (paginator.Page[tagentity.Tag], error) {
	return u.tagRepo.GetPage(ctx, pageSize, pageNum)
}

func (u *TagUseCase) GetTagsByPostID(ctx context.Context, postID int64) ([]tagentity.Tag, error) {
	return u.tagRepo.GetTagsByPostID(ctx, postID)
}
