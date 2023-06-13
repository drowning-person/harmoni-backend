package service

import (
	"context"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/usecase"

	"go.uber.org/zap"
)

type PostService struct {
	pc     *usecase.PostUseCase
	logger *zap.SugaredLogger
}

func NewPostService(
	pc *usecase.PostUseCase,
	logger *zap.SugaredLogger,
) *PostService {
	return &PostService{
		pc:     pc,
		logger: logger,
	}
}

func (s *PostService) GetPosts(ctx context.Context, req *postentity.GetPostsRequest) (*postentity.GetPostsReply, error) {
	posts, err := s.pc.GetPage(ctx, req.PageSize, req.Page, req.Order)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	res := paginator.Page[postentity.PostDetail]{
		CurrentPage: posts.CurrentPage,
		PageSize:    posts.PageSize,
		Total:       posts.Total,
		Pages:       posts.Pages,
		Data:        make([]postentity.PostDetail, 0, len(posts.Data)),
	}

	for _, post := range posts.Data {
		res.Data = append(res.Data, postentity.ConvertPostToDisplayDetail(&post))
	}

	return &postentity.GetPostsReply{
		Page: res,
	}, nil
}

func (s *PostService) GetPostDetail(ctx context.Context, req *postentity.GetPostDetailRequest) (*postentity.GetPostDetailReply, error) {
	post, exist, err := s.pc.GetByPostID(ctx, req.PostID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	} else if !exist {
		return nil, errorx.NotFound(reason.PostNotFound)
	}

	return &postentity.GetPostDetailReply{
		PostDetail: postentity.ConvertPostToDisplayDetail(post),
	}, nil
}

func (s *PostService) Create(ctx context.Context, req *postentity.CreatePostRequest) (*postentity.CreatePostReply, error) {
	post := postentity.Post{
		AuthorID: req.UserID,
		TagID:    req.TagID,
		Title:    req.Title,
		Content:  req.Content,
	}

	post, err := s.pc.Create(ctx, &post)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &postentity.CreatePostReply{
		PostDetail: postentity.ConvertPostToDisplayDetail(&post),
	}, nil
}
