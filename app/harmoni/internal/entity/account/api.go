package account

import emailentity "harmoni/app/harmoni/internal/entity/email"

type AccountActionType uint8

const (
	RegisterAct AccountActionType = iota + 1
	ChangeEmailAct
	ChangePasswordAct
	ResetPasswordAct
	BindEmailAct
)

func AccountActionTypeToEmailType(act AccountActionType) emailentity.EmailType {
	switch act {
	case RegisterAct, BindEmailAct:
		return emailentity.BindEmail
	case ChangeEmailAct:
		return emailentity.ChangeEmail
	case ChangePasswordAct:
		return emailentity.ChangePassword
	case ResetPasswordAct:
		return emailentity.ResetPassword
	}

	return 0
}

type MailSendRequest struct {
	Email  string            `json:"email" validate:"required,email"`
	UserID int64             `json:"-"`
	Type   AccountActionType `json:"type" validate:"required"`
}

type MailSendReply struct {
}

type MailCheckRequest struct {
	Email  string            `json:"email" validate:"required,email"`
	UserID int64             `json:"-"`
	Type   AccountActionType `json:"type,omitempty" validate:"required"`
	Code   string            `json:"code,omitempty" validate:"required"`
}

type MailCheckReply struct {
}

type ChangeEmailRequest struct {
	UserID   int64  `json:"user_id,string,omitempty" validate:"required"`
	NewEmail string `json:"new_email,omitempty" validate:"required,email"`
	Code     string `json:"code,omitempty" validate:"required"`
}

type ChangeEmailReply struct {
}

type ChangePasswordRequest struct {
	UserID      int64  `json:"user_id,string,omitempty" validate:"required"`
	OldPassword string `json:"old_password,omitempty" validate:"required" label:"旧密码"`
	NewPassword string `json:"new_password,omitempty" validate:"required" label:"新密码"`
}

type ChangePasswordReply struct {
}

type ResetPasswordRequest struct {
	UserID      int64  `json:"user_id,string,omitempty" validate:"required"`
	NewPassword string `json:"new_password,omitempty" validate:"required"`
}

type ResetPasswordReply struct {
}
