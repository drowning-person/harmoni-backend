package like

import (
	pb "harmoni/api/like/grpc/v1"
	entitylike "harmoni/app/like/internal/entity/like"
)

func convertDomainToReply(like *entitylike.Like) *pb.LikeEntity {
	return &pb.LikeEntity{
		ObjectType:   like.ObjectType,
		ObjectID:     like.ObjectID,
		TargetUserID: like.TargetUser.Id,
	}
}
