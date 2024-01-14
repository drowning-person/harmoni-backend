package remind

import (
	"context"
	"errors"
	v1 "harmoni/app/harmoni/api/grpc/v1/user"
	"harmoni/app/notification/internal/entity/remind"
	"harmoni/app/notification/internal/infrastructure/po/notification"
	"harmoni/internal/pkg/data"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/paginator"
	"harmoni/internal/pkg/reason"
	"harmoni/internal/pkg/set"
	"harmoni/internal/types/action"
	"harmoni/internal/types/iface"
	"harmoni/internal/types/object"
	"harmoni/internal/types/persistence"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
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
	err := r.data.DB(ctx).
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
			err = r.data.DB(ctx).Create(&remind).Error
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
	err = r.data.DB(ctx).Select("sender_id").
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
		err = r.data.DB(ctx).Model(&notification.RemindParticipant{}).
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
		err = r.data.DB(ctx).Create(participants).Error
		if err != nil {
			return errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}
	return nil
}

func (r *RemindRepo) ListNRemindParticipant(ctx context.Context, remindIDs []int64, n int) ([]*notification.RemindParticipant, error) {
	participants := []*notification.RemindParticipant{}
	var subQuery = r.data.DB(ctx).
		Select("COUNT(*)").
		Table("remind_participant").
		Where("remind_id = rp.remind_id").
		Where("rp_id <= rp.rp_id")
	err := r.data.DB(ctx).Table("remind_participant AS rp").
		Where("(?) <= ?", subQuery, n).
		Where("rp.remind_id in ?", remindIDs).Order("rp.rp_id").
		Find(&participants).Error
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return participants, nil
}

func (r *RemindRepo) List(ctx context.Context, req *remind.ListReq) (*paginator.Page[*remind.Remind], error) {
	page := paginator.NewPage[*notification.NotifyRemind](req.Page, req.Size)
	err := page.SelectPages(r.data.DB(ctx).Scopes(
		WithRecipientID(req.UserID),
		WithAction(req.Action)).
		Order("updated_at"))
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	remindIDs := make([]int64, page.Total)
	for i := range page.Data {
		remindIDs[i] = page.Data[i].RemindID
	}
	senders, err := r.ListNRemindParticipant(ctx, remindIDs, req.SenderCount)
	if err != nil {
		return nil, err
	}
	reminds := buildDomain(page.Data)
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
	return &paginator.Page[*remind.Remind]{
		CurrentPage: page.CurrentPage,
		PageSize:    page.PageSize,
		Total:       page.Total,
		Data:        reminds,
	}, nil
}

func (r *RemindRepo) UpdateLastReadTime(ctx context.Context, req *remind.UpdateLastReadTimeReq) error {
	err := r.data.DB(ctx).Model(&notification.NotifyRemind{}).
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
	db := r.data.DB(ctx)
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

type remindSenderResult struct {
	SenderID  int64     `gorm:"column:sender_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (r *RemindRepo) ListRemindSenders(ctx context.Context, req *remind.ListRemindSendersReq) (*paginator.Page[*remind.RemindSender], error) {
	result := paginator.NewPageFromReq[*remindSenderResult](req.PageRequest)
	db := r.data.DB(ctx).Table("remind_participant AS rp").
		Select("rp.sender_id", "rp.created_at").
		Where("rp.remind_id = ?", req.RemindID).
		Where("rp.action = ?", req.Action).
		Order("rp.rp_id").
		Scan(result)
	err := result.SelectPages(db)
	if err != nil {
		return nil, errorx.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	resp, err := r.uc.List(ctx, &v1.ListBasicsRequest{Ids: lo.Map(result.Data,
		func(r *remindSenderResult, _ int) int64 {
			return r.SenderID
		})})
	if err != nil {
		return nil, errorx.InternalServer(reason.ServerError).WithError(err).WithStack()
	}
	senders := v1.UserBasicList(resp.GetUsers()).ToMap()
	items := paginator.NewPageFromReq[*remind.RemindSender](req.PageRequest)
	datas := make([]*remind.RemindSender, len(result.Data))
	for i := range datas {
		datas[i].Sender = senders[result.Data[i].SenderID]
		datas[i].CreatedAt = &result.Data[i].CreatedAt
	}
	return items, nil
}
