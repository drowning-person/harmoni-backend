package service

import (
	"context"
	"harmoni/internal/entity/paginator"
	tagentity "harmoni/internal/entity/tag"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type TagService struct {
	tc     *usecase.TagUseCase
	logger *zap.SugaredLogger
}

func NewTagService(
	tc *usecase.TagUseCase,
	logger *zap.SugaredLogger,
) *TagService {
	return &TagService{
		tc:     tc,
		logger: logger,
	}
}

// GetTags TODO: Add condition
func (s *TagService) GetTags(ctx context.Context, req *tagentity.GetTagsRequest) (tagentity.GetTagsReply, error) {
	tags, err := s.tc.GetPage(ctx, req.PageSize, req.Page)
	if err != nil {
		s.logger.Errorln(err)
		return tagentity.GetTagsReply{}, err
	}

	res := paginator.Page[tagentity.TagInfo]{
		CurrentPage: tags.CurrentPage,
		PageSize:    tags.PageSize,
		Total:       tags.Total,
		Pages:       tags.Pages,
		Data:        make([]tagentity.TagInfo, 0, len(tags.Data)),
	}

	for _, tag := range tags.Data {
		res.Data = append(res.Data, tagentity.ConvertTagToDisplay(tag))
	}

	return tagentity.GetTagsReply{
		Page: res,
	}, nil
}

func (s *TagService) Create(ctx context.Context, req *tagentity.CreateTagRequest) (tagentity.CreateTagReply, error) {
	exist, err := s.tc.IsTagExistByName(ctx, req.Name)
	if err != nil {
		s.logger.Errorln(err)
		return tagentity.CreateTagReply{}, err
	}
	if exist {
		s.logger.Infof("Create Tag attempt failed. Tag with name '%v' already exists.\n", req.Name)
		return tagentity.CreateTagReply{}, errorx.BadRequest(reason.TagAlreadyExist)
	}

	tag := tagentity.Tag{
		TagName:      req.Name,
		Introduction: req.Introduction,
	}
	tag, err = s.tc.Create(ctx, &tag)
	if err != nil {
		s.logger.Errorln(err)
		return tagentity.CreateTagReply{}, err
	}

	return tagentity.CreateTagReply{
		TagDetail: tagentity.ConvertTagToDetailDisplay(tag),
	}, nil
}

func (s *TagService) GetTag(ctx context.Context, req *tagentity.GetTagDetailRequest) (tagentity.GetTagDetailReply, error) {
	tag, err := s.tc.GetTagByTagID(ctx, req.TagID)
	if err != nil {
		s.logger.Errorln(err)
		return tagentity.GetTagDetailReply{}, err
	}

	return tagentity.GetTagDetailReply{
		TagDetail: tagentity.ConvertTagToDetailDisplay(tag),
	}, nil
}
