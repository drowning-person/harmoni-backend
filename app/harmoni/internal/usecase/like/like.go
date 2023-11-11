package like

import (
	"context"
	commententity "harmoni/app/harmoni/internal/entity/comment"
	likeentity "harmoni/app/harmoni/internal/entity/like"
	"harmoni/app/harmoni/internal/entity/paginator"
	postentity "harmoni/app/harmoni/internal/entity/post"
	userentity "harmoni/app/harmoni/internal/entity/user"
	"harmoni/app/harmoni/internal/infrastructure/config"
	"harmoni/app/harmoni/internal/pkg/reason"
	event "harmoni/app/harmoni/internal/types/events/like"
	"harmoni/app/harmoni/internal/types/iface"
	"harmoni/app/harmoni/internal/usecase/like/events"
	postuse "harmoni/app/harmoni/internal/usecase/post"
	"harmoni/internal/pkg/errorx"

	"time"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSetLikeUsecase = wire.NewSet(
	NewLikeUsecase,
	events.NewLikeEventsHandler,
)

type LikeUsecase struct {
	likeRepo    likeentity.LikeRepository
	userRepo    userentity.UserRepository
	commentRepo commententity.CommentRepository
	postUseCase *postuse.PostUseCase
	logger      *zap.SugaredLogger

	publisher iface.Publisher
}

func NewLikeUsecase(
	conf *config.MessageQueue,
	likeRepo likeentity.LikeRepository,
	postUseCase *postuse.PostUseCase,
	commentRepo commententity.CommentRepository,
	userRepo userentity.UserRepository,
	logger *zap.SugaredLogger,
	publisher iface.Publisher,
) (*LikeUsecase, func(), error) {
	lc := &LikeUsecase{
		likeRepo:    likeRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		postUseCase: postUseCase,
		logger:      logger.With("module", "usecase/like"),
		publisher:   publisher,
	}

	return lc, func() {
		lc.publisher.Close()
	}, nil
}

func (u *LikeUsecase) Like(ctx context.Context, like *likeentity.Like, isCancel bool) error {
	var (
		exist        bool
		err          error
		post         *postentity.PostInfo
		comment      *commententity.Comment
		targetUserID int64
	)
	// TODO: cache post and comment existence
	switch like.LikeType {
	case likeentity.LikePost:
		post, exist, err = u.postUseCase.GetByPostID(ctx, 0, like.LikingID)
		targetUserID = post.User.UserID
	case likeentity.LikeComment:
		comment, exist, err = u.commentRepo.GetByCommentID(ctx, like.LikingID)
		targetUserID = comment.Author.UserID
	default:
		return errorx.BadRequest(reason.LikeUnknownType)
	}
	if err != nil {
		return err
	} else if !exist {
		return errorx.NotFound(reason.ObjectNotFound)
	}

	like.TargetUserID = targetUserID
	err = u.likeRepo.Like(ctx, like, like.TargetUserID, isCancel)
	if err != nil {
		return err
	}

	now := time.Now()
	msg := event.LikeCreatedMessage{
		BaseMessage: event.BaseMessage{
			LikeType: like.LikeType.ToEventLikeType(),
		},
		TargetUserID: like.TargetUserID,
		UserID:       like.UserID,
		IsCancel:     isCancel,
		LikingID:     like.LikingID,
		CreatedAt:    &now,
	}

	err = u.publisher.Publish(ctx, event.TopicLikeCreated, &msg)
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err)
	}

	return nil
}

func (u *LikeUsecase) ListLikingIDs(ctx context.Context, query *likeentity.LikeQuery) (paginator.Page[int64], error) {
	if query.Type == likeentity.LikeComment {
		return paginator.Page[int64]{}, errorx.BadRequest(reason.TypeNotSupport)
	}
	return u.likeRepo.ListLikingIDs(ctx, query)
}

func (u *LikeUsecase) IsLiking(ctx context.Context, like *likeentity.Like) (bool, error) {
	return u.likeRepo.IsLiking(ctx, like)
}

func (u *LikeUsecase) LikeCount(ctx context.Context, like *likeentity.Like) (int64, error) {
	count, exist, err := u.likeRepo.LikeCount(ctx, like)
	if err != nil {
		return 0, err
	} else if exist {
		return count, nil
	}

	switch like.LikeType {
	case likeentity.LikePost:
		count, exist, err = u.commentRepo.GetLikeCount(ctx, like.LikingID)
	case likeentity.LikeComment:
		count, exist, err = u.postUseCase.GetLikeCount(ctx, like.LikingID)
	case likeentity.LikeUser:
		count, exist, err = u.userRepo.GetLikeCount(ctx, like.LikingID)
	}
	if err != nil {
		return 0, err
	}
	if !exist {
		count = -1
	}

	err = u.likeRepo.CacheLikeCount(ctx, like, count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *LikeUsecase) BatchLikeCount(ctx context.Context, likeType likeentity.LikeType) (map[int64]int64, error) {
	return u.likeRepo.BatchLikeCount(ctx, likeType)
}

func (u *LikeUsecase) BatchLikeCountByIDs(ctx context.Context, likingIDs []int64, likeType likeentity.LikeType) (map[int64]int64, error) {
	return u.likeRepo.BatchLikeCountByIDs(ctx, likingIDs, likeType)
}

func (u *LikeUsecase) GetLikingObjects(ctx context.Context, userID int64, objectIDs []int64, likeType likeentity.LikeType) (any, error) {
	switch likeType {
	case likeentity.LikePost:
		postInfos, err := u.postUseCase.BatchByIDs(ctx, objectIDs)
		if err != nil {
			return nil, err
		}

		postInfos, err = u.postUseCase.MergeList(ctx, userID, postInfos)
		if err != nil {
			return nil, err
		}

		return postInfos, nil
	default:
		return nil, errorx.BadRequest(reason.TypeNotSupport)
	}
}
