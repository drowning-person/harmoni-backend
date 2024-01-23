package like

import (
	"context"

	objectv1 "harmoni/api/common/object/v1"
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
		ObjectType:     req.GetObjectType(),
		TargetObjectID: req.GetObjectID(),
		IsCancel:       req.IsCancel,
	})
	if err != nil {
		s.logger.Error(err)
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
		UserID:     req.GetUserID(),
		ObjectType: req.GetObjectType(),
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	reply := pb.LikeListReply{
		PageRely: paginator.NewPageReply(
			req.GetPageRequest().GetNum(),
			req.GetPageRequest().GetSize(), total),
		LikeList: make([]*pb.LikeEntity, len(list)),
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
		s.logger.Error(err)
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

func (s *LikeService) LikeCount(ctx context.Context, req *objectv1.Object) (*pb.LikeCountReply, error) {
	count, err := s.lu.LikeCount(ctx, req)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &pb.LikeCountReply{
		Count: count.Count,
	}, nil
}

func (s *LikeService) ListLikeCount(ctx context.Context, req *pb.ListLikeCountRequest) (*pb.ListLikeCountReply, error) {
	counts, err := s.lu.ListLikeCount(ctx, req.GetObjectIds(), req.GetObjectType())
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}
	return &pb.ListLikeCountReply{Counts: counts.ToMap()}, nil
}
