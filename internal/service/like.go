package service

import (
	"context"
	likeentity "harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type LikeService struct {
	lc     *usecase.LikeUsecase
	logger *zap.SugaredLogger
}

func NewLikeUsecase(
	lc *usecase.LikeUsecase,
	logger *zap.SugaredLogger) *LikeService {
	return &LikeService{
		lc:     lc,
		logger: logger.With("module", "service/like"),
	}
}
func (s *LikeService) Like(ctx context.Context, req *likeentity.LikeRequest) (*likeentity.LikeReply, error) {
	err := s.lc.Like(ctx, &likeentity.Like{
		UserID:   req.UserID,
		LikingID: req.ObjectID,
		LikeType: req.Type,
	}, req.IsCancel)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &likeentity.LikeReply{}, nil
}

func (s *LikeService) GetLikings(ctx context.Context, req *likeentity.GetLikingsRequest) (*likeentity.GetLikingsReply[any], error) {
	idPage, err := s.lc.ListLikingIDs(ctx, &likeentity.LikeQuery{
		PageCond: req.PageCond,
		UserID:   req.UserID,
		Type:     req.Type,
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	res := paginator.Page[any]{
		CurrentPage: idPage.CurrentPage,
		PageSize:    idPage.PageSize,
		Total:       idPage.Total,
		Pages:       idPage.Pages,
	}

	res.Data, err = s.lc.GetLikingObjects(ctx, idPage.Data, req.Type)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &likeentity.GetLikingsReply[any]{Page: res}, nil
}

func (s *LikeService) IsLiking(ctx context.Context, req *likeentity.IsLikingRequest) (*likeentity.IsLikingReply, error) {
	liking, err := s.lc.IsLiking(ctx, &likeentity.Like{
		UserID:   req.UserID,
		LikingID: req.ObjectID,
		LikeType: req.Type,
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &likeentity.IsLikingReply{Liking: liking}, nil
}
