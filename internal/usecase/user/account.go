package user

import (
	"context"
	accountentity "harmoni/internal/entity/account"
	"harmoni/internal/entity/auth"
	emailentity "harmoni/internal/entity/email"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/infrastructure/config"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/usecase/email"
	"time"

	"go.uber.org/zap"
)

type AccountUsecase struct {
	conf         *config.Email
	authUsecase  *AuthUseCase
	userRepo     userentity.UserRepository
	userUseCase  *UserUseCase
	emailUsecase *email.EmailUsecase
	logger       *zap.Logger
}

func NewAccountUsecase(
	conf *config.Email,
	authUsecase *AuthUseCase,
	userRepo userentity.UserRepository,
	emailUsecase *email.EmailUsecase,
	userUseCase *UserUseCase,
	logger *zap.Logger) *AccountUsecase {
	return &AccountUsecase{
		conf:         conf,
		authUsecase:  authUsecase,
		userRepo:     userRepo,
		emailUsecase: emailUsecase,
		userUseCase:  userUseCase,
		logger:       logger,
	}
}

// modifyAccount modifies the user account based on the specified action type.
func (u *AccountUsecase) modifyAccountByEmail(ctx context.Context, olduser *userentity.User, newUser *userentity.User, actionType accountentity.AccountActionType) error {
	user, exist, err := u.userRepo.GetByUserID(ctx, olduser.UserID)
	if err != nil {
		return err
	} else if !exist {
		return errorx.BadRequest(reason.UserNotFound)
	}

	// Check if the user's email is verified before performing the action.
	status, err := u.userRepo.GetModifyStaus(ctx, user.UserID, userentity.VerifyByEmail, actionType)
	if err != nil {
		return err
	} else if status != userentity.VerifiedEmail {
		return errorx.BadRequest(reason.EmailNeedToBeVerifiedBeforeAct)
	}

	switch actionType {
	case accountentity.ChangeEmailAct:
		user.Email = newUser.Email
		return u.userRepo.ModifyEmail(ctx, user)
	case accountentity.ChangePasswordAct, accountentity.ResetPasswordAct:
		if actionType == accountentity.ChangePasswordAct && !u.userUseCase.CheckPassword(user.Password, olduser.Password) {
			return errorx.BadRequest(reason.OldPasswordVerificationFailed)
		}
		if olduser.Password == newUser.Password {
			return errorx.BadRequest(reason.NewPasswordSameAsPreviousSetting)
		}
		user.Password, err = u.userUseCase.HashAndSalt(newUser.Password)
		if err != nil {
			return err
		}

		err = u.userRepo.ModifyPassword(ctx, user)
		if err != nil {
			return err
		}

		err = u.authUsecase.RevokeTokens(ctx, user.UserID, auth.RefreshTokenType)
		if err != nil {
			u.logger.Sugar().Errorf("revoke refresh tokens failed %s", err)
		}

		err = u.authUsecase.RevokeTokens(ctx, user.UserID, auth.AccessTokenType)
		if err != nil {
			u.logger.Sugar().Errorf("revoke access tokens failed %s", err)
		}

	}

	return nil
}

func (u *AccountUsecase) ChangeEmailByEmail(ctx context.Context, user *userentity.User, newUser *userentity.User) error {
	return u.modifyAccountByEmail(ctx, user, newUser, accountentity.ChangeEmailAct)
}

func (u *AccountUsecase) ChangePasswordByEmail(ctx context.Context, user *userentity.User, newUser *userentity.User) error {
	return u.modifyAccountByEmail(ctx, user, newUser, accountentity.ChangePasswordAct)
}

func (u *AccountUsecase) ResetPasswordByEmail(ctx context.Context, user *userentity.User, newUser *userentity.User) error {
	return u.modifyAccountByEmail(ctx, user, newUser, accountentity.ResetPasswordAct)
}

func (u *AccountUsecase) CheckVerificationCodeByEmail(ctx context.Context, user *userentity.User, code string, actionType accountentity.AccountActionType) error {
	emailType := accountentity.AccountActionTypeToEmailType(actionType)

	err := u.emailUsecase.VerifyCode(ctx, user.Email, code, emailType)
	if err != nil {
		return err
	}

	if emailType == emailentity.BindEmail {
		return nil
	}

	err = u.userRepo.SetModifyStatus(ctx, user.UserID, userentity.VerifiedEmail, userentity.VerifyByEmail, actionType, u.conf.CodeTTL)
	if err != nil {
		return err
	}

	return nil
}

func (u *AccountUsecase) SendVerificationCodeByEmail(ctx context.Context, user *userentity.User, actionType accountentity.AccountActionType) error {
	err := u.emailUsecase.CheckBeforeSendCode(ctx, user.Email, accountentity.AccountActionTypeToEmailType(actionType))
	if err != nil {
		return err
	}

	_, exist, err := u.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	switch actionType {
	case accountentity.ChangeEmailAct, accountentity.ChangePasswordAct, accountentity.ResetPasswordAct:
		if !exist {
			return errorx.BadRequest(reason.UserNotFound)
		}
	case accountentity.RegisterAct, accountentity.BindEmailAct:
		if exist {
			return errorx.BadRequest(reason.EmailDuplicate)
		}
	}

	code := u.emailUsecase.GenCode(ctx)

	var title, body string
	switch actionType {
	case accountentity.RegisterAct:
		title, body, err = u.emailUsecase.RegisterTemplate(ctx, code)
		if err != nil {
			return err
		}
	case accountentity.ChangeEmailAct, accountentity.ChangePasswordAct:
		changeName := "密码"
		if actionType == accountentity.ChangeEmailAct {
			changeName = "邮箱"
		}
		title, body, err = u.emailUsecase.ChangeTemplate(ctx, code, changeName)
		if err != nil {
			return err
		}
	case accountentity.ResetPasswordAct:
		title, body, err = u.emailUsecase.ResetPasswordTemplate(ctx, code)
		if err != nil {
			return err
		}
	}

	data := emailentity.EmailCodeContent{
		Code:        code,
		LastReqTime: time.Now().Unix(),
	}

	err = u.emailUsecase.SendAndSaveCode(ctx, user.Email, title, body, data.ToJSONString(), accountentity.AccountActionTypeToEmailType(actionType))
	if err != nil {
		return err
	}

	return nil
}
