package usecase

import (
	"context"
	"encoding/json"
	"harmoni/internal/conf"
	"harmoni/internal/entity"
	commententity "harmoni/internal/entity/comment"
	likeentity "harmoni/internal/entity/like"
	"harmoni/internal/entity/paginator"
	postentity "harmoni/internal/entity/post"
	userentity "harmoni/internal/entity/user"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/queue/rabbitmq"
	"harmoni/internal/pkg/reason"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	jsonType = "application/json"
)

type LikeUsecase struct {
	likeRepo    likeentity.LikeRepository
	postRepo    postentity.PostRepository
	userRepo    userentity.UserRepository
	commentRepo commententity.CommentRepository
	consumers   *rabbitmq.RabbitListener
	producers   *rabbitmq.RabbitMqSender
	logger      *zap.SugaredLogger
}

func NewLikeUsecase(
	conf *conf.MessageQueue,
	likeRepo likeentity.LikeRepository,
	postRepo postentity.PostRepository,
	commentRepo commententity.CommentRepository,
	userRepo userentity.UserRepository,
	logger *zap.SugaredLogger) (*LikeUsecase, func(), error) {
	logger.Debugf("conf is %#v, rabbitmq conf is %#v", *conf, conf.RabbitMQ)
	lc := &LikeUsecase{
		likeRepo:    likeRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		userRepo:    userRepo,
		logger:      logger.With("module", "usecase/like"),
	}
	if conf.RabbitMQ != nil {
		rabbitConf := rabbitmq.RabbitConf{
			Username: conf.RabbitMQ.Username,
			Password: conf.RabbitMQ.Password,
			Host:     conf.RabbitMQ.Host,
			Port:     conf.RabbitMQ.Port,
			VHost:    conf.RabbitMQ.VHost,
		}

		admin := rabbitmq.MustNewAdmin(rabbitConf)
		err := admin.DeclareExchange(rabbitmq.ExchangeConf{
			ExchangeName: entity.LikeExchange,
			Type:         amqp.ExchangeDirect,
			Durable:      true,
		}, nil)
		if err != nil {
			return nil, func() {}, err
		}
		err = admin.DeclareQueue(rabbitmq.QueueConf{
			Name:    entity.LikeQueue,
			Durable: true,
		}, nil)
		if err != nil {
			return nil, func() {}, err
		}
		err = admin.Bind(entity.LikeQueue, entity.LikeBindKey, entity.LikeExchange, false, nil)
		if err != nil {
			return nil, func() {}, err
		}

		lc.consumers = rabbitmq.MustNewListener(rabbitmq.RabbitListenerConf{
			RabbitConf: rabbitConf,
			ListenerQueues: []rabbitmq.ConsumerConf{
				{
					Name:      entity.LikeQueue,
					AutoAck:   false,
					Exclusive: false,
					NoLocal:   false,
					NoWait:    false,
				},
			},
		}, lc)

		go lc.consumers.Start()

		lc.producers = rabbitmq.MustNewSender(rabbitmq.RabbitSenderConf{
			RabbitConf:  rabbitConf,
			ContentType: jsonType,
		})

	}

	return lc, func() { lc.consumers.Stop() }, nil
}

func (u *LikeUsecase) Like(ctx context.Context, like *likeentity.Like, isCancel bool) error {
	var (
		exist        bool
		err          error
		post         *postentity.Post
		comment      *commententity.Comment
		targetUserID int64
	)
	// TODO: cache post and comment existence
	switch like.LikeType {
	case likeentity.LikePost:
		post, exist, err = u.postRepo.GetByPostID(ctx, like.LikingID)
		targetUserID = post.AuthorID
	case likeentity.LikeComment:
		comment, exist, err = u.commentRepo.GetByCommentID(ctx, like.LikingID)
		targetUserID = comment.AuthorID
	}
	if err != nil {
		return err
	} else if !exist {
		return errorx.NotFound(reason.ObjectNotFound)
	}

	err = u.likeRepo.Like(ctx, like, targetUserID, isCancel)
	if err != nil {
		return err
	}

	now := time.Now()
	msg := likeentity.LikeMessage{
		Type:     likeentity.LikeActionMessage,
		LikeType: like.LikeType,
		ActionMessage: &likeentity.ActionMessage{
			UserID:    like.UserID,
			IsCancel:  isCancel,
			LikingID:  like.LikingID,
			CreatedAt: &now,
		},
	}
	data, _ := json.Marshal(msg)
	err = u.producers.Send(ctx, entity.LikeExchange, entity.LikeBindKey, data)
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
		count, exist, err = u.postRepo.GetLikeCount(ctx, like.LikingID)
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

func (u *LikeUsecase) GetLikingObjects(ctx context.Context, objectIDs []int64, likeType likeentity.LikeType) ([]any, error) {
	switch likeType {
	case likeentity.LikePost:
		posts, err := u.postRepo.BatchBasicInfoByIDs(ctx, objectIDs)
		if err != nil {
			return nil, err
		}

		likeCounts, err := u.likeRepo.BatchLikeCountByIDs(ctx, objectIDs, likeType)
		if err != nil {
			return nil, err
		}

		objects := make([]any, len(posts))
		for i := 0; i < len(posts); i++ {
			posts[i].LikeCount = likeCounts[posts[i].PostID]
			objects[i] = postentity.ConvertPostToDisplay(&posts[i])
		}
		return objects, nil
	default:
		return []any{}, errorx.BadRequest(reason.TypeNotSupport)
	}
}

func (u *LikeUsecase) Consume(message []byte) error {
	msg := likeentity.LikeMessage{}
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	u.logger.Debugf("receive msg %#v", msg)
	ctx := context.Background()
	switch msg.Type {
	case likeentity.LikeCountMessage:
		switch msg.LikeType {
		case likeentity.LikePost:
			for likingID, likeCount := range msg.CountMessage.Counts {
				err = u.postRepo.UpdateLikeCount(ctx, likingID, likeCount)
			}
		case likeentity.LikeComment:
			for likingID, likeCount := range msg.CountMessage.Counts {
				err = u.commentRepo.UpdateLikeCount(ctx, likingID, likeCount)
			}
		case likeentity.LikeUser:
			for likingID, likeCount := range msg.CountMessage.Counts {
				err = u.userRepo.UpdateLikeCount(ctx, likingID, likeCount)
			}
		}
	case likeentity.LikeActionMessage:
		err = u.likeRepo.Save(ctx, &likeentity.Like{
			UserID:   msg.ActionMessage.UserID,
			LikingID: msg.ActionMessage.LikingID,
			LikeType: msg.LikeType,
		}, msg.ActionMessage.IsCancel)
	}

	return err
}
