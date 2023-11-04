package service

import (
	"context"
	accountentity "harmoni/internal/entity/account"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/usecase/user"

	"go.uber.org/zap"
)

type AccountService struct {
	ac     *user.AccountUsecase
	logger *zap.SugaredLogger
}

func NewAccountService(ac *user.AccountUsecase, logger *zap.SugaredLogger) *AccountService {
	return &AccountService{
		ac:     ac,
		logger: logger,
	}
}

func (s *AccountService) MailSend(ctx context.Context, req *accountentity.MailSendRequest) (*accountentity.MailSendReply, error) {
	err := s.ac.SendVerificationCodeByEmail(ctx, &userentity.User{UserID: req.UserID, Email: req.Email}, req.Type)
	if err != nil {
		return nil, err
	}

	return &accountentity.MailSendReply{}, nil
}

func (s *AccountService) MailCheck(ctx context.Context, req *accountentity.MailCheckRequest) (*accountentity.MailCheckReply, error) {
	err := s.ac.CheckVerificationCodeByEmail(ctx, &userentity.User{UserID: req.UserID, Email: req.Email}, req.Code, req.Type)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	return &accountentity.MailCheckReply{}, nil
}

func (s *AccountService) ChangeEmail(ctx context.Context, req *accountentity.ChangeEmailRequest) (*accountentity.ChangeEmailReply, error) {
	user := &userentity.User{UserID: req.UserID, Email: req.NewEmail}
	err := s.ac.CheckVerificationCodeByEmail(ctx, user, req.Code, accountentity.BindEmailAct)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	err = s.ac.ChangeEmailByEmail(ctx, user, &userentity.User{Email: req.NewEmail})
	if err != nil {
		return nil, err
	}

	return &accountentity.ChangeEmailReply{}, nil
}

func (s *AccountService) ChangePassword(ctx context.Context, req *accountentity.ChangePasswordRequest) (*accountentity.ChangePasswordReply, error) {
	err := s.ac.ChangePasswordByEmail(ctx, &userentity.User{UserID: req.UserID, Password: req.OldPassword}, &userentity.User{Password: req.NewPassword})
	if err != nil {
		return nil, err
	}

	return &accountentity.ChangePasswordReply{}, nil
}

func (s *AccountService) ResetPassword(ctx context.Context, req *accountentity.ResetPasswordRequest) (*accountentity.ResetPasswordReply, error) {
	err := s.ac.ResetPasswordByEmail(ctx, &userentity.User{UserID: req.UserID}, &userentity.User{Password: req.NewPassword})
	if err != nil {
		return nil, err
	}

	return &accountentity.ResetPasswordReply{}, nil
}
