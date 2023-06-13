package usecase

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	"harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"

	"go.uber.org/zap"
)

type CommentUseCase struct {
	commentRepo commententity.CommentRepository
	likeUsecase *LikeUsecase
	logger      *zap.SugaredLogger
}

func NewCommentUseCase(
	commentRepo commententity.CommentRepository,
	likeUsecase *LikeUsecase,
	logger *zap.SugaredLogger) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		likeUsecase: likeUsecase,
		logger:      logger,
	}
}

func (u *CommentUseCase) Create(ctx context.Context, comment *commententity.Comment) error {
	return u.commentRepo.Create(ctx, comment)
}

func (u *CommentUseCase) GetPage(ctx context.Context, commentQuery *commententity.CommentQuery) (paginator.Page[commententity.Comment], error) {
	comments, err := u.commentRepo.GetPage(ctx, commentQuery)
	if err != nil {
		return paginator.Page[commententity.Comment]{}, err
	}

	commentIDs := make([]int64, len(comments.Data))
	for i, comment := range comments.Data {
		commentIDs[i] = comment.CommentID
	}
	likes, err := u.likeUsecase.BatchLikeCountByIDs(ctx, commentIDs, like.LikeComment)
	if err != nil {
		return paginator.Page[commententity.Comment]{}, err
	}
	for i := range comments.Data {
		comments.Data[i].LikeCount = likes[comments.Data[i].CommentID]
	}

	return comments, err
}
