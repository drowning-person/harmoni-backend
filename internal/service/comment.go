package service

import (
	"context"
	commententity "harmoni/internal/entity/comment"
	"harmoni/internal/entity/paginator"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type CommentService struct {
	cc     *usecase.CommentUseCase
	logger *zap.SugaredLogger
}

func NewCommentService(
	cc *usecase.CommentUseCase,
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

	res := paginator.Page[commententity.CommentDetail]{
		CurrentPage: comments.CurrentPage,
		PageSize:    comments.PageSize,
		Total:       comments.Total,
		Pages:       comments.Pages,
		Data:        make([]commententity.CommentDetail, 0, len(comments.Data)),
	}

	for _, comment := range comments.Data {
		res.Data = append(res.Data, commententity.CommentToCommentDetail(&comment))
	}

	return commententity.GetCommentsReply{
		Page: res,
	}, nil
}

func (s *CommentService) Create(ctx context.Context, req *commententity.CreateCommentRequest) (commententity.CreateCommentReply, error) {
	comment := commententity.Comment{
		AuthorID: req.UserID,
		ParentID: req.ParentID,
		RootID:   req.RootID,
		Content:  req.Content,
	}

	err := s.cc.Create(ctx, &comment)
	if err != nil {
		s.logger.Errorln(err)
		return commententity.CreateCommentReply{}, err
	}

	return commententity.CreateCommentReply{
		CommentDetail: commententity.CommentToCommentDetail(&comment),
	}, nil
}
