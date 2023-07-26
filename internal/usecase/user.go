package usecase

import (
	"context"
	"harmoni/internal/entity/like"
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
	fileUsecase *FileUseCase
	likeUsecase *LikeUsecase
	logger      *zap.SugaredLogger
	reg         *regexp.Regexp
}

func NewUserUseCase(
	userRepo userentity.UserRepository,
	authUsecase *AuthUseCase,
	fileUsecase *FileUseCase,
	likeUsecase *LikeUsecase,
	logger *zap.SugaredLogger,
) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		authUsecase: authUsecase,
		fileUsecase: fileUsecase,
		likeUsecase: likeUsecase,
		logger:      logger,
		reg:         regexp.MustCompile("^[-_!a-zA-Z0-9\u4e00-\u9fa5]+$"),
	}
}

func (u *UserUseCase) HashAndSalt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func (u *UserUseCase) CheckPassword(userPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	return err == nil
}

func (u *UserUseCase) Create(ctx context.Context, user *userentity.User) error {
	var err error
	user.Password, err = u.HashAndSalt(user.Password)
	if err != nil {
		return err
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return err
}

func (u *UserUseCase) VerifyPassword(ctx context.Context, password, hashedPwd string) error {
	if !u.CheckPassword(hashedPwd, password) {
		return errorx.BadRequest(reason.EmailOrPasswordWrong)
	}
	return nil
}

func (u *UserUseCase) GetByUserID(ctx context.Context, userID int64) (*userentity.User, bool, error) {
	user, exist, err := u.userRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	count, err := u.likeUsecase.LikeCount(ctx, &like.Like{LikingID: userID, LikeType: like.LikeUser})
	if err != nil {
		return nil, false, err
	}
	user.LikeCount = count

	return user, exist, nil
}

func (u *UserUseCase) GetByUserIDs(ctx context.Context, userIDs []int64) ([]userentity.User, error) {
	return u.userRepo.GetByUserIDs(ctx, userIDs)
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*userentity.User, bool, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func (u *UserUseCase) GetPage(ctx context.Context, pageSize int64, pageNum int64) (paginator.Page[userentity.User], error) {
	return u.userRepo.GetPage(ctx, pageSize, pageNum)
}

func (u *UserUseCase) SetAvatar(ctx context.Context, userID int64, fileID int64) error {
	err := u.userRepo.SetAvatarID(ctx, userID, fileID)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetAvatarLink(ctx context.Context, userID int64) (string, error) {
	avatarID, err := u.userRepo.GetAvatarID(ctx, userID)
	if err != nil {
		if errx, ok := err.(*errorx.Error); ok && errorx.IsNotFound(errx) {
			return u.fileUsecase.GetFileLink(ctx, 0)
		}
		return "", err
	}
	return u.fileUsecase.GetFileLink(ctx, avatarID)
}

// ValidUsername 验证用户
func (u *UserUseCase) ValidUsername(username string) error {
	if !u.reg.MatchString(username) || utf8.RuneCountInString(username) > 32 {
		return errorx.BadRequest(reason.UsernameInvalid)
	}
	return nil
}
