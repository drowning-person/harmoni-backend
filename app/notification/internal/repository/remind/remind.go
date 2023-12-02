package remind

import (
	"context"
	"errors"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/app/notification/internal/entity/remind"
	"harmoni/app/notification/internal/infrastructure/po/notification"
	"harmoni/internal/pkg/data"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/set"
	"harmoni/internal/types/action"
	"harmoni/internal/types/iface"
	"harmoni/internal/types/object"
	"harmoni/internal/types/persistence"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ remind.RemindRepository = (*RemindRepo)(nil)

type RemindRepo struct {
	uc         v1.UserClient
	data       *data.DB
	uniqueRepo iface.UniqueIDRepository
	logger     *log.Helper
}

func New(
	uc v1.UserClient,
	data *data.DB,
	uniqueRepo iface.UniqueIDRepository,
	logger log.Logger,
) (*RemindRepo, error) {
	r := &RemindRepo{
		data:       data,
		uniqueRepo: uniqueRepo,
		logger:     log.NewHelper(log.With(logger, "module", "repository/remind")),
	}
	r.uc = uc
	return r, nil
}

func fromCreate(notifyRemind *notification.NotifyRemind, req *remind.CreateReq) {
	notifyRemind.Action = int8(req.Action)
	notifyRemind.RecipientID = req.RecipientID
	notifyRemind.ObjectID = req.ObjectID
	notifyRemind.ObjectType = int8(req.ObjectType)
	notifyRemind.Content = req.Content
	if req.LastReadTime != nil {
		notifyRemind.LastReadTime = *req.LastReadTime
	}
}

func buildDomain(notifyReminds []*notification.NotifyRemind) []*remind.Remind {
	reminds := make([]*remind.Remind, 0, len(notifyReminds))
	for i := range notifyReminds {
		reminds[i] = &remind.Remind{
			RemindID:     notifyReminds[i].RemindID,
			Recipient:    &v1.UserBasic{Id: notifyReminds[i].RecipientID},
			Action:       action.Action(notifyReminds[i].Action),
			ObjectID:     notifyReminds[i].ObjectID,
			ObjectType:   object.ObjectType(notifyReminds[i].ObjectType),
			Content:      notifyReminds[i].Content,
			LastReadTime: notifyReminds[i].LastReadTime,
		}
	}
	return reminds
}

func WithRecipientID(recipientID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("recipient_id = ?", recipientID)
	}
}

func WithAction(act action.Action) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if act == action.ActionNo {
			return db
		}
		return db.Where("action = ?", act)
	}
}

func (r *RemindRepo) Create(ctx context.Context, req *remind.CreateReq) error {
	participants := make([]*notification.RemindParticipant, len(req.SenderIDs))
	remind := notification.NotifyRemind{}
	err := r.data.DB(ctx).WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Scopes(WithRecipientID(req.RecipientID),
			WithAction(req.Action),
			persistence.ByObject(req.ObjectID, req.ObjectType)).
		First(&remind).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fromCreate(&remind, req)
			remind.RemindID, err = r.uniqueRepo.GenUniqueID(ctx)
			if err != nil {
				return err
			}
			err = r.data.DB(ctx).WithContext(ctx).Create(&remind).Error
			if err != nil {
				return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
		} else {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}

	// if sender not stored, add to remind_participant
	// otherwise, update time
	storedSenders := make([]int64, 0)
	err = r.data.DB(ctx).WithContext(ctx).
		Select("sender_id").
		Table("remind_participant").
		Where("remind_id = ?", remind.RemindID).
		Find(&storedSenders).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	notStoredSenders := make([]int64, 0)
	storedSendersMap := make(map[int64]bool)
	for i := range storedSenders {
		storedSendersMap[storedSenders[i]] = true
	}
	for i := range req.SenderIDs {
		if !storedSendersMap[req.SenderIDs[i]] {
			notStoredSenders = append(notStoredSenders, req.SenderIDs[i])
		}
	}
	if len(notStoredSenders) != len(req.SenderIDs) {
		needUpdateSenders := make([]int64, 0)
		for i := range req.SenderIDs {
			if storedSendersMap[req.SenderIDs[i]] {
				needUpdateSenders = append(needUpdateSenders, req.SenderIDs[i])
			}
		}
		updatedTime := time.Now()
		if req.LastReadTime != nil {
			updatedTime = *req.LastReadTime
		}
		err = r.data.DB(ctx).WithContext(ctx).
			Model(&notification.RemindParticipant{}).
			Where("sender_id IN (?)", needUpdateSenders).
			UpdateColumn("updated_at", updatedTime).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}

	if len(notStoredSenders) > 0 {
		for i := range notStoredSenders {
			uniqueID, err := r.uniqueRepo.GenUniqueID(ctx)
			if err != nil {
				return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
			participants[i] = &notification.RemindParticipant{
				SenderID: req.SenderIDs[i],
				RemindID: remind.RemindID,
				RpID:     uniqueID,
			}
		}
		err = r.data.DB(ctx).WithContext(ctx).Create(participants).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}
	return nil
}

func (r *RemindRepo) List(ctx context.Context, req *remind.ListReq) ([]*remind.Remind, error) {
	notifyReminds := []*notification.NotifyRemind{}
	err := r.data.DB(ctx).WithContext(ctx).Scopes(
		WithRecipientID(req.UserID),
		WithAction(req.Action)).
		Order("updated_at").
		Find(&notifyReminds).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	remindIDs := make([]int64, len(notifyReminds))
	for i := range notifyReminds {
		remindIDs[i] = notifyReminds[i].RemindID
	}
	senders := []*notification.RemindParticipant{}
	err = r.data.DB(ctx).WithContext(ctx).Where("remind_id IN (?)", remindIDs).
		Find(&senders).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	reminds := buildDomain(notifyReminds)
	userIDSet := set.New[int64]()
	for i := range senders {
		userIDSet.Add(senders[i].SenderID)
	}
	userIDSet.Add(req.UserID)
	users, err := r.uc.List(ctx, &v1.ListBasicsRequest{Ids: userIDSet.ToArray()})
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	userMap := v1.UserBasicList(users.GetUsers()).ToMap()
	sendersMap := map[int64][]*v1.UserBasic{}
	for i := range senders {
		sendersMap[senders[i].RemindID] = append(sendersMap[senders[i].RemindID], userMap[senders[i].SenderID])
	}
	for i := range reminds {
		reminds[i].Recipient = userMap[reminds[i].Recipient.GetId()]
		reminds[i].Senders = sendersMap[reminds[i].RemindID]
	}
	return reminds, nil
}

func (r *RemindRepo) UpdateLastReadTime(ctx context.Context, req *remind.UpdateLastReadTimeReq) error {
	err := r.data.DB(ctx).WithContext(ctx).Model(&notification.NotifyRemind{}).
		Scopes(
			WithRecipientID(req.UserID),
			WithAction(req.Action),
		).
		UpdateColumn("last_read_time", req.ReadTime).Error
	if err != nil {
		return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *RemindRepo) Count(ctx context.Context, req *remind.CountReq) (int64, error) {
	var count int64
	db := r.data.DB(ctx).WithContext(ctx)
	if req.UnRead {
		db = db.Joins("JOIN remind_participant AS rp ON notify_remind.remind_id = rp.remind_id AND rp.created_at > notify_remind.last_read_time")
	}
	err := db.Model(&notification.NotifyRemind{}).
		Scopes(
			WithRecipientID(req.UserID),
			WithAction(req.Action),
		).Count(&count).Error
	if err != nil {
		return 0, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}
