package service

import (
	"context"
	followentity "harmoni/internal/entity/follow"
	"harmoni/internal/entity/paginator"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type FollowService struct {
	fc     *usecase.FollowUseCase
	uc     *usecase.UserUseCase
	logger *zap.SugaredLogger
}

func NewFollowService(fc *usecase.FollowUseCase, uc *usecase.UserUseCase, logger *zap.SugaredLogger) *FollowService {
	return &FollowService{
		fc:     fc,
		uc:     uc,
		logger: logger,
	}
}

func (s *FollowService) Follow(ctx context.Context, req *followentity.FollowRequest) (*followentity.FollowReply, error) {
	var err error
	if req.IsCancel {
		err = s.fc.FollowCancel(ctx, &followentity.Follow{
			FollowerID:   req.UserID,
			FollowingID:  req.ObjectID,
			FollowedType: req.Type,
		})
	} else {
		err = s.fc.Follow(ctx, &followentity.Follow{
			FollowerID:   req.UserID,
			FollowingID:  req.ObjectID,
			FollowedType: req.Type,
		})
	}
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &followentity.FollowReply{}, nil
}

func (s *FollowService) GetFollowers(ctx context.Context, req *followentity.GetFollowersRequest) (*followentity.GetFollowerReply, error) {
	userIDPage, err := s.fc.GetFollowers(ctx, &followentity.FollowQuery{
		PageCond: req.PageCond,
		UserID:   req.UserID,
		Type:     followentity.FollowUser,
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	res := paginator.Page[userentity.UserBasicInfo]{
		CurrentPage: userIDPage.CurrentPage,
		PageSize:    userIDPage.PageSize,
		Total:       userIDPage.Total,
		Pages:       userIDPage.Pages,
		Data:        make([]userentity.UserBasicInfo, len(userIDPage.Data)),
	}

	users, err := s.uc.GetByUserIDs(ctx, userIDPage.Data)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	for i, v := range users {
		link, err := s.uc.GetAvatarLink(ctx, v.UserID)
		if err != nil {
			s.logger.Errorln(err)
			return nil, err
		}
		res.Data[i] = v.ToBasicInfo(link)
	}

	return &followentity.GetFollowerReply{Page: res}, nil
}

func (s *FollowService) GetFollowing(ctx context.Context, req *followentity.GetFollowingsRequest) (*followentity.GetFollowingsReply[any], error) {
	idPage, err := s.fc.GetFollowings(ctx, &followentity.FollowQuery{
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

	res.Data, err = s.fc.GetFollowingObjects(ctx, idPage.Data, req.Type)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &followentity.GetFollowingsReply[any]{Page: res}, nil
}

func (s *FollowService) IsFollowing(ctx context.Context, req *followentity.IsFollowingRequest) (*followentity.IsFollowingReply, error) {
	following, err := s.fc.IsFollowing(ctx, &followentity.Follow{
		FollowerID:   req.UserID,
		FollowingID:  req.ObjectID,
		FollowedType: req.Type,
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &followentity.IsFollowingReply{Following: following}, nil
}

func (s *FollowService) AreFollowEachOther(ctx context.Context, req *followentity.AreFollowEachOtherRequest) (*followentity.AreFollowEachOtherReply, error) {
	following, err := s.fc.AreFollowEachOther(ctx, req.UserID, req.ObjectID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &followentity.AreFollowEachOtherReply{FollowedEachOther: following}, nil
}
