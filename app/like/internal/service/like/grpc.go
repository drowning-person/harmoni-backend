package like

import (
	"context"

	pb "harmoni/api/like/grpc/v1"
	entitylike "harmoni/app/like/internal/entity/like"
	"harmoni/app/like/internal/usecase/like"
	"harmoni/internal/pkg/paginator"

	"github.com/go-kratos/kratos/v2/log"
)

type LikeService struct {
	pb.UnimplementedLikeServer

	lu     *like.LikeUsecase
	logger *log.Helper
}

func NewLikeService(
	lu *like.LikeUsecase,
	logger log.Logger,
) *LikeService {
	return &LikeService{
		lu:     lu,
		logger: log.NewHelper(log.With(logger, "module", "service/like", "service", "like")),
	}
}

func (s *LikeService) Like(ctx context.Context, req *pb.LikeRequest) (*pb.LikeReply, error) {
	err := s.lu.Like(ctx, &like.LikeRequest{
		UserID:         req.GetUserID(),
		TargetUserID:   req.GetTargetUserID(),
		LikeType:       entitylike.LikeType(req.GetLikeType()),
		TargetObjectID: req.GetObjectID(),
		IsCancel:       req.IsCancel,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LikeReply{}, nil
}

func (s *LikeService) UserLikeList(ctx context.Context, req *pb.UserLikeListRequest) (*pb.LikeListReply, error) {
	list, total, err := s.lu.ListLikeObjectByUserID(ctx, &entitylike.ListLikeObjectQuery{
		PageRequest: paginator.PageRequest{
			Num:  req.GetPageRequest().GetNum(),
			Size: req.GetPageRequest().GetSize(),
		},
		UserID: req.GetUserID(),
	})
	if err != nil {
		return nil, err
	}

	size := int64(len(list))
	reply := pb.LikeListReply{
		PageRely: paginator.NewPageReply(
			req.GetPageRequest().GetNum(),
			size, total),
		LikeList: make([]*pb.LikeEntity, size),
	}
	for i := range list {
		reply.LikeList[i] = convertDomainToReply(list[i])
	}

	return &reply, nil
}

func (s *LikeService) ObjectLikeList(ctx context.Context, req *pb.ObjectLikeListRequest) (*pb.LikeListReply, error) {
	list, total, err := s.lu.ListLikeUserByObjectID(ctx,
		&entitylike.ListObjectLikedUserQuery{
			PageRequest: paginator.PageRequest{
				Num:  req.GetPageRequest().GetNum(),
				Size: req.GetPageRequest().GetSize(),
			},
			ObjectID: req.GetObjectID(),
		})
	if err != nil {
		return nil, err
	}

	size := int64(len(list))
	reply := pb.LikeListReply{
		PageRely: paginator.NewPageReply(
			req.GetPageRequest().GetNum(),
			size, total),
		LikeList: make([]*pb.LikeEntity, size),
	}
	for i := range list {
		reply.LikeList[i] = convertDomainToReply(list[i])
	}
	return &reply, nil
}
