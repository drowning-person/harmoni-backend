package usecase

import (
	"context"
	"harmoni/internal/entity/paginator"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"regexp"
	"unicode/utf8"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo    userentity.UserRepository
	authUsecase *AuthUseCase
	logger      *zap.SugaredLogger
	reg         *regexp.Regexp
}

func NewUserUseCase(userRepo userentity.UserRepository, authUsecase *AuthUseCase, logger *zap.SugaredLogger) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		authUsecase: authUsecase,
		logger:      logger,
		reg:         regexp.MustCompile("^[-_!a-zA-Z0-9\u4e00-\u9fa5]+$"),
	}
}

func hashAndSalt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func checkPassword(userPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	return err == nil
}

func (u *UserUseCase) Create(ctx context.Context, user *userentity.User) (userentity.User, error) {
	var err error
	user.Password, err = hashAndSalt(user.Password)
	if err != nil {
		return userentity.User{}, err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return userentity.User{}, err
	}

	return *user, err
}

func (u *UserUseCase) VerifyPassword(ctx context.Context, password, hashedPwd string) error {
	if !checkPassword(hashedPwd, password) {
		return errorx.BadRequest(reason.EmailOrPasswordWrong)
	}
	return nil
}

func (u *UserUseCase) GetByUserID(ctx context.Context, userID int64) (*userentity.User, bool, error) {
	return u.userRepo.GetByUserID(ctx, userID)
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*userentity.User, bool, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func (u *UserUseCase) GetPage(ctx context.Context, pageSize int64, pageNum int64) (paginator.Page[userentity.User], error) {
	return u.userRepo.GetPage(ctx, pageSize, pageNum)
}

// ValidUsername 验证用户
func (u *UserUseCase) ValidUsername(username string) error {
	if !u.reg.MatchString(username) || utf8.RuneCountInString(username) > 32 {
		return errorx.BadRequest(reason.UsernameInvalid)
	}
	return nil
}
