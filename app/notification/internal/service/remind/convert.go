package remind

import (
	v1 "harmoni/app/notification/api/http/v1/notification"
	"harmoni/app/notification/internal/entity/remind"
)

func ConverRemindToResp(remind *remind.Remind) *v1.Remind {
	return &v1.Remind{
		RemindID:   remind.RemindID,
		Recipient:  remind.Recipient,
		ObjectID:   remind.ObjectID,
		ObjectType: remind.ObjectType,
		Content:    remind.Content,
	}
}
