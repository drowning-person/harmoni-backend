package service

import (
	"context"
	accountentity "harmoni/internal/entity/account"
	authentity "harmoni/internal/entity/auth"
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
	auc    *usecase.AccountUsecase
	logger *zap.SugaredLogger
}

func NewUserService(
	uc *usecase.UserUseCase,
	ac *usecase.AuthUseCase,
	auc *usecase.AccountUsecase,
	logger *zap.SugaredLogger,
) *UserService {
	return &UserService{
		uc:     uc,
		ac:     ac,
		auc:    auc,
		logger: logger,
	}
}

func (s *UserService) GetUserByUserID(ctx context.Context, req *userentity.GetUserDetailRequest) (*userentity.GetUserDetailReply, error) {
	user, exist, err := s.uc.GetByUserID(ctx, req.UserID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	} else if !exist {
		return nil, errorx.NotFound(reason.UserNotFound)
	}
	link, err := s.uc.GetAvatarLink(ctx, user.UserID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &userentity.GetUserDetailReply{
		UserDetail: userentity.ConvertUserToDetailDisplay(user, link),
	}, nil
}

// GetUsers TODO: Add condition
func (s *UserService) GetUsers(ctx context.Context, req *userentity.GetUsersRequest) (*userentity.GetUsersReply, error) {
	users, err := s.uc.GetPage(ctx, req.PageSize, req.Page)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	res := paginator.Page[userentity.UserBasicInfo]{
		CurrentPage: users.CurrentPage,
		PageSize:    users.PageSize,
		Total:       users.Total,
		Pages:       users.Pages,
		Data:        make([]userentity.UserBasicInfo, 0, len(users.Data)),
	}

	for _, user := range users.Data {
		link, err := s.uc.GetAvatarLink(ctx, user.UserID)
		if err != nil {
			s.logger.Errorln(err)
			return nil, err
		}
		res.Data = append(res.Data, user.ToBasicInfo(link))
	}

	return &userentity.GetUsersReply{
		Page: res,
	}, nil
}

func (s *UserService) RegisterByEmail(ctx context.Context, req *userentity.UserRegisterRequest) (*userentity.UserRegisterReply, error) {
	_, exist, err := s.uc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	} else if exist {
		s.logger.Infof("Registration attempt failed. User with email '%v' already exists.\n", req.Email)
		return nil, errorx.BadRequest(reason.EmailDuplicate)
	}

	if err := s.uc.ValidUsername(req.Name); err != nil {
		return nil, err
	}

	user := userentity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	err = s.auc.CheckVerificationCodeByEmail(ctx, &user, req.RegisterCode, accountentity.RegisterAct)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = s.uc.Create(ctx, &user)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	token, err := s.ac.GenToken(ctx, &user, authentity.AccessTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	refreshToken, err := s.ac.GenToken(ctx, &user, authentity.RefreshTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	link, err := s.uc.GetAvatarLink(ctx, user.UserID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &userentity.UserRegisterReply{
		User:         user.ToBasicInfo(link),
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *userentity.UserLoginRequset) (*userentity.UserLoginReply, error) {
	user, exist, err := s.uc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	} else if !exist {
		s.logger.Infof("Login attempt failed. User with email '%v' not found.\n", req.Email)
		return nil, errorx.NotFound(reason.UserNotFound)
	}

	err = s.uc.VerifyPassword(ctx, req.Password, user.Password)
	if err != nil {
		s.logger.Infof("Login attempt failed. User (%#v) password is wrong.\n", user)
		return nil, err
	}

	token, err := s.ac.GenToken(ctx, user, authentity.AccessTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	refreshToken, err := s.ac.GenToken(ctx, user, authentity.RefreshTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	link, err := s.uc.GetAvatarLink(ctx, user.UserID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &userentity.UserLoginReply{
		User:         user.ToBasicInfo(link),
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, req *userentity.UserLogoutRequest) (*userentity.UserLogoutReply, error) {
	err := s.ac.RevokeToken(ctx, req.UserID, req.AccessTokenID, authentity.AccessTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	claims, err := s.ac.VerifyToken(ctx, req.RefreshToken, authentity.RefreshTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = s.ac.RevokeToken(ctx, req.UserID, claims.ID, authentity.RefreshTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userentity.UserLogoutReply{}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, req *userentity.RefreshTokenRequest) (*userentity.RefreshTokenReply, error) {
	claims, err := s.ac.VerifyToken(ctx, req.RefreshToken, authentity.RefreshTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	accessToken, err := s.ac.GenToken(ctx, &userentity.User{UserID: claims.UserID, Name: claims.Name}, authentity.AccessTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	newRefreshToken, err := s.ac.ReGenToken(ctx, claims, authentity.RefreshTokenType)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &userentity.RefreshTokenReply{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
