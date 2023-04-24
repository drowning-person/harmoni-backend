package service

import (
	"context"
	"harmoni/internal/entity/paginator"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type UserService struct {
	uc     *usecase.UserUseCase
	ac     *usecase.AuthUseCase
	logger *zap.SugaredLogger
}

func NewUserService(
	uc *usecase.UserUseCase,
	ac *usecase.AuthUseCase,
	logger *zap.SugaredLogger,
) *UserService {
	return &UserService{
		uc:     uc,
		ac:     ac,
		logger: logger,
	}
}

func (s *UserService) GetUserByUserID(ctx context.Context, userID int64) (userentity.UserDetail, error) {
	user, err := s.uc.GetUserByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorln(err)
		return userentity.UserDetail{}, err
	}

	return userentity.ConvertUserToDetailDisplay(user), nil
}

// GetUsers TODO: Add condition
func (s *UserService) GetUsers(ctx context.Context, pageSize, pageNum int64) (paginator.Page[userentity.BasicUserInfo], error) {
	users, err := s.uc.GetPage(ctx, pageSize, pageNum)
	if err != nil {
		s.logger.Errorln(err)
		return paginator.Page[userentity.BasicUserInfo]{}, err
	}

	res := paginator.Page[userentity.BasicUserInfo]{
		CurrentPage: users.CurrentPage,
		PageSize:    users.PageSize,
		Total:       users.Total,
		Pages:       users.Pages,
		Data:        make([]userentity.BasicUserInfo, 0, len(users.Data)),
	}

	for _, user := range users.Data {
		res.Data = append(res.Data, userentity.ConvertUserToDisplay(user))
	}

	return res, nil
}

func (s *UserService) Register(ctx context.Context, req *userentity.UserRegisterRequest) (userentity.UserRegisterReply, error) {
	exist, err := s.uc.IsUserExistByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorln(err)
		return userentity.UserRegisterReply{}, err
	}
	if exist {
		s.logger.Infof("Registration attempt failed. User with email '%v' already exists.\n", req.Email)
		return userentity.UserRegisterReply{}, errorx.BadRequest(reason.EmailDuplicate)
	}

	if err := s.uc.ValidUsername(req.Name); err != nil {
		return userentity.UserRegisterReply{}, err
	}

	user := userentity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	user, err = s.uc.Create(ctx, &user)
	if err != nil {
		s.logger.Errorln(err)
		return userentity.UserRegisterReply{}, err
	}

	token, err := s.ac.GenToken(ctx, user.UserID, user.Name)
	if err != nil {
		s.logger.Errorln(err)
		return userentity.UserRegisterReply{}, err
	}

	return userentity.UserRegisterReply{
		User:        userentity.ConvertUserToDisplay(user),
		AccessToken: token,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *userentity.UserLoginRequset) (userentity.UserLoginReply, error) {
	user, err := s.uc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Errorln(err)
		return userentity.UserLoginReply{}, err
	}

	err = s.uc.VerifyPassword(ctx, req.Password, user.Password)
	if err != nil {
		s.logger.Infof("Login attempt failed. User (%#v) password is wrong.\n", user)
		return userentity.UserLoginReply{}, err
	}

	token, err := s.ac.GenToken(ctx, user.UserID, user.Name)
	if err != nil {
		s.logger.Errorln(err)
		return userentity.UserLoginReply{}, err
	}

	return userentity.UserLoginReply{
		User:        userentity.ConvertUserToDisplay(user),
		AccessToken: token,
	}, nil
}
