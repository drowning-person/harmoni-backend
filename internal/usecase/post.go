package usecase

import (
	"context"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"

	"go.uber.org/zap"
)

type PostUseCase struct {
	postRepo postentity.PostRepository
	logger   *zap.SugaredLogger
}

func NewPostUseCase(postRepo postentity.PostRepository, logger *zap.SugaredLogger) *PostUseCase {
	return &PostUseCase{
		postRepo: postRepo,
		logger:   logger,
	}
}

func (u *PostUseCase) Create(ctx context.Context, post *postentity.Post) (postentity.Post, error) {
	err := u.postRepo.Create(ctx, post)
	if err != nil {
		return postentity.Post{}, err
	}

	return *post, err
}

func (u *PostUseCase) GetByPostID(ctx context.Context, postID int64) (*postentity.Post, bool, error) {
	return u.postRepo.GetByPostID(ctx, postID)
}

func (u *PostUseCase) GetBasicInfoByPostID(ctx context.Context, postID int64) (*postentity.Post, bool, error) {
	return u.postRepo.GetBasicInfoByPostID(ctx, postID)
}

// func (u *PostUseCase) BatchByIDs(ctx context.Context, postIDs []int64) ([]Post, error)
func (u *PostUseCase) GetPage(ctx context.Context, pageSize, pageNum int64, orderCond string) (paginator.Page[postentity.Post], error) {
	posts, err := u.postRepo.GetPage(ctx, pageSize, pageNum, orderCond)
	if err != nil {
		return paginator.Page[postentity.Post]{}, err
	}

	return posts, err
}

func (u *PostUseCase) LikePost(ctx context.Context, postID int64, userID int64, direction int8) error {
	err := u.postRepo.LikePost(ctx, postID, userID, direction)
	if err != nil {
		return err
	}

	return err
}
