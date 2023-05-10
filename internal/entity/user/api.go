package user

import "harmoni/internal/entity/paginator"

type GetUserDetailRequest struct {
	UserID int64 `params:"id" validate:"required"`
}

type GetUserDetailReply struct {
	UserDetail
}

type GetUsersRequest struct {
	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
	// query condition
	QueryCond string `query:"cond" validate:"omitempty,oneof=newest" label:"排序"`
}

type GetUsersReply struct {
	paginator.Page[BasicUserInfo]
}

type UserLoginRequset struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginReply struct {
	User         BasicUserInfo
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserLogoutRequest struct {
	UserID        int64  `json:"user_id,omitempty"`
	AccessTokenID string `json:"access_token_id,omitempty"`
	RefreshToken  string `json:"refresh_token,omitempty" validate:"required"`
}

type UserLogoutReply struct {
}

type UserRegisterRequest struct {
	Name         string `json:"name" validate:"required,gte=3,lte=20"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required"`
	RegisterCode string `json:"register_code" validate:"required"`
}

type UserRegisterReply struct {
	User         BasicUserInfo
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type UserResetPasswordRequest struct {
	UserID int64 `json:"user_id,omitempty"`
}

type UserResetPasswordReply struct {
}

type UserSendCodeByEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserSendCodeByEmailReply struct {
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token,omitempty" validate:"required"`
}

type RefreshTokenReply struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
