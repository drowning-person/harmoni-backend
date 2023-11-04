package service

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	"harmoni/internal/usecase/comment"

	"go.uber.org/zap"
)

type CommentService struct {
	cc     *comment.CommentUseCase
	logger *zap.SugaredLogger
}

func NewCommentService(
	cc *comment.CommentUseCase,
	logger *zap.SugaredLogger,
) *CommentService {
	return &CommentService{
		cc:     cc,
		logger: logger,
	}
}

func (s *CommentService) GetComments(ctx context.Context, req *commententity.GetCommentsRequest) (commententity.GetCommentsReply, error) {
	q := commententity.ConvertPageReqToCommentQuery(req)
	comments, err := s.cc.GetPage(ctx, &q)
	if err != nil {
		s.logger.Errorln(err)
		return commententity.GetCommentsReply{}, err
	}

	return commententity.GetCommentsReply{
		Page: *comments,
	}, nil
}

func (s *CommentService) Create(ctx context.Context, req *commententity.CreateCommentRequest) (commententity.CreateCommentReply, error) {
	comment := req.ToDomain()

	err := s.cc.Create(ctx, comment)
	if err != nil {
		s.logger.Errorln(err)
		return commententity.CreateCommentReply{}, err
	}

	return commententity.CreateCommentReply{
		Comment: *comment,
	}, nil
}
