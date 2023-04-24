package usecase

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	"harmoni/internal/entity/paginator"

	"go.uber.org/zap"
)

type CommentUseCase struct {
	commentRepo commententity.CommentRepository
	logger      *zap.SugaredLogger
}

func NewCommentUseCase(commentRepo commententity.CommentRepository, logger *zap.SugaredLogger) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (u *CommentUseCase) Create(ctx context.Context, comment *commententity.Comment) (commententity.Comment, error) {
	err := u.commentRepo.Create(ctx, comment)
	if err != nil {
		return commententity.Comment{}, err
	}

	return *comment, nil
}

func (u *CommentUseCase) GetPage(ctx context.Context, commentQuery *commententity.CommentQuery) (paginator.Page[commententity.Comment], error) {
	comments, err := u.commentRepo.GetPage(ctx, commentQuery)
	if err != nil {
		return paginator.Page[commententity.Comment]{}, err
	}

	return comments, err
}
