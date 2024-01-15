package like

import (
	"context"

	pb "harmoni/api/like/grpc/v1"
	entitylike "harmoni/app/like/internal/entity/like"
	"harmoni/app/like/internal/usecase/like"

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
