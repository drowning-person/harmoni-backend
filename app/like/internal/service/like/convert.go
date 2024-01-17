package like

import (
	pb "harmoni/api/like/grpc/v1"
	v1 "harmoni/api/like/mq/v1"
	entitylike "harmoni/app/like/internal/entity/like"
)

func convertDomainToReply(like *entitylike.Like) *pb.LikeEntity {
	return &pb.LikeEntity{
		LikeType:     v1.LikeType(like.LikeType),
		ObjectID:     like.ObjectID,
		TargetUserID: like.TargetUser.Id,
	}
}
