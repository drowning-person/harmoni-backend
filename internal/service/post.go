package service

import (
	"context"
	postentity "harmoni/internal/entity/post"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/usecase/post"
	"harmoni/internal/usecase/tag"

	"go.uber.org/zap"
)

type PostService struct {
	pc     *post.PostUseCase
	tc     *tag.TagUseCase
	logger *zap.SugaredLogger
}

func NewPostService(
	pc *post.PostUseCase,
	tc *tag.TagUseCase,
	logger *zap.SugaredLogger,
) *PostService {
	return &PostService{
		pc:     pc,
		tc:     tc,
		logger: logger,
	}
}

func (s *PostService) GetPosts(ctx context.Context, req *postentity.GetPostsRequest) (*postentity.GetPostsReply, error) {
	posts, err := s.pc.GetPage(ctx, &postentity.PostQuery{PageCond: req.PageCond, QueryCond: req.Order, TagID: req.TagID, UserID: req.UserID})
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &postentity.GetPostsReply{
		Page: posts,
	}, nil
}

func (s *PostService) GetPostInfo(ctx context.Context, req *postentity.GetPostInfoRequest) (*postentity.GetPostInfoReply, error) {
	post, exist, err := s.pc.GetByPostID(ctx, req.UserID, req.PostID)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	} else if !exist {
		return nil, errorx.NotFound(reason.PostNotFound)
	}

	return &postentity.GetPostInfoReply{
		PostInfo: *post,
	}, nil
}

func (s *PostService) Create(ctx context.Context, req *postentity.CreatePostRequest) (*postentity.CreatePostReply, error) {
	post := postentity.Post{
		AuthorID: req.UserID,
		TagIDs:   req.TagIDs,
		Title:    req.Title,
		Content:  req.Content,
	}

	postInfo, err := s.pc.Create(ctx, &post)
	if err != nil {
		s.logger.Errorln(err)
		return nil, err
	}

	return &postentity.CreatePostReply{
		PostInfo: *postInfo,
	}, nil
}
