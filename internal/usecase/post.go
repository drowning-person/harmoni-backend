package usecase

import (
	"context"
	"harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"

	"go.uber.org/zap"
)

type PostUseCase struct {
	postRepo    postentity.PostRepository
	likeUsecase *LikeUsecase
	logger      *zap.SugaredLogger
}

func NewPostUseCase(postRepo postentity.PostRepository, likeUsecase *LikeUsecase, logger *zap.SugaredLogger) *PostUseCase {
	return &PostUseCase{
		postRepo:    postRepo,
		likeUsecase: likeUsecase,
		logger:      logger,
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
	post, exist, err := u.postRepo.GetByPostID(ctx, postID)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, false, nil
	}
	post.LikeCount, err = u.likeUsecase.LikeCount(ctx, &like.Like{LikingID: postID, LikeType: like.LikePost})
	if err != nil {
		return nil, false, err
	}

	return post, true, nil
}

func (u *PostUseCase) GetBasicInfoByPostID(ctx context.Context, postID int64) (*postentity.Post, bool, error) {
	return u.postRepo.GetBasicInfoByPostID(ctx, postID)
}

// func (u *PostUseCase) BatchByIDs(ctx context.Context, postIDs []int64) ([]Post, error)
func (u *PostUseCase) GetPage(ctx context.Context, queryCond *postentity.PostQuery) (paginator.Page[postentity.Post], error) {
	posts, err := u.postRepo.GetPage(ctx, queryCond)
	if err != nil {
		return paginator.Page[postentity.Post]{}, err
	}
	postIDs := make([]int64, len(posts.Data))
	for i, post := range posts.Data {
		postIDs[i] = post.PostID
	}
	likes, err := u.likeUsecase.BatchLikeCountByIDs(ctx, postIDs, like.LikePost)
	if err != nil {
		return paginator.Page[postentity.Post]{}, err
	}
	for i := range posts.Data {
		posts.Data[i].LikeCount = likes[posts.Data[i].PostID]
	}

	return posts, err
}

func (u *PostUseCase) BatchBasicInfoByIDs(ctx context.Context, postIDs []int64) ([]postentity.Post, error) {
	return u.postRepo.BatchBasicInfoByIDs(ctx, postIDs)
}

/* func (u *PostUseCase) LikePost(ctx context.Context, postID int64, userID int64, direction int8) error {
	err := u.postRepo.LikePost(ctx, postID, userID, direction)
	if err != nil {
		return err
	}

	return err
}
*/
