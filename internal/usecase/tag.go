package usecase

import (
	"context"
	"errors"
	"harmoni/internal/entity/paginator"
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/pkg/errorx"

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

func (u *TagUseCase) Create(ctx context.Context, tag *tagentity.Tag) (tagentity.Tag, error) {
	err := u.tagRepo.Create(ctx, tag)
	if err != nil {
		return tagentity.Tag{}, err
	}

	return *tag, err
}

func (u *TagUseCase) GetTagByTagID(ctx context.Context, tagID int64) (tagentity.Tag, error) {
	tag, err := u.tagRepo.GetByTagID(ctx, tagID)
	if err != nil {
		return tagentity.Tag{}, err
	}
	return tag, nil
}

func (u *TagUseCase) IsTagExistByName(ctx context.Context, name string) (bool, error) {
	tag, err := u.tagRepo.GetByTagName(ctx, name)
	if err != nil {
		myErr := &errorx.Error{}
		if errors.As(err, &myErr) {
			if errorx.IsNotFound(err.(*errorx.Error)) {
				return false, nil
			}
		}
		return false, err
	}
	if tag.ID == 0 {
		return false, nil
	}

	return true, nil
}

func (u *TagUseCase) GetPage(ctx context.Context, pageSize int64, pageNum int64) (paginator.Page[tagentity.Tag], error) {
	return u.tagRepo.GetPage(ctx, pageSize, pageNum)
}
