package user

import (
	"context"

	pb "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/app/harmoni/internal/entity/auth"
	"harmoni/app/harmoni/internal/pkg/reason"
	"harmoni/app/harmoni/internal/usecase/file"
	"harmoni/app/harmoni/internal/usecase/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/set"
)

type UserGRPCService struct {
	pb.UnimplementedUserServer
	ac *user.AuthUseCase
	fc *file.FileUseCase
	uc *user.UserUseCase
}

func NewUserGRPCService(
	ac *user.AuthUseCase,
	fc *file.FileUseCase,
	uc *user.UserUseCase,
) *UserGRPCService {
	return &UserGRPCService{
		ac: ac,
		fc: fc,
		uc: uc,
	}
}

func (s *UserGRPCService) GetBasic(ctx context.Context, req *pb.GetBasicRequest) (*pb.UserBasic, error) {
	user, exist, err := s.uc.GetByUserID(ctx, req.GetId())
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, errorx.NotFound(reason.UserNotFound)
	}
	link, err := s.fc.GetFileLink(ctx, user.Avatar)
	if err != nil {
		return nil, err
	}
	basic := user.ToBasicInfo(link)
	return (&pb.UserBasic{}).FromDomain(&basic), nil
}

func (s *UserGRPCService) List(ctx context.Context, req *pb.ListBasicsRequest) (*pb.ListBasicsResponse, error) {
	users, err := s.uc.GetByUserIDs(ctx, req.GetIds())
	if err != nil {
		return nil, err
	}
	avatarSet := set.New[int64]()
	for i := range users {
		avatarSet.Add(users[i].Avatar)
	}
	avatarlinkMap, err := s.fc.ListFileLinkMap(ctx, avatarSet.ToArray())
	if err != nil {
		return nil, err
	}
	return &pb.ListBasicsResponse{
		Users: pb.ListFromDomain(users.ToUserBasics(avatarlinkMap)),
	}, nil
}

func (s *UserGRPCService) VerifyToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	claims, err := s.ac.VerifyToken(ctx, req.GetToken(), auth.AccessTokenType)
	if err != nil {
		return nil, err
	}
	return &pb.TokenResponse{
		User: &pb.UserBasic{
			Id:   claims.TokenInfo.UserID,
			Name: claims.TokenInfo.Name,
		},
	}, nil
}
