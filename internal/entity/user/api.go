package user

import "harmoni/internal/entity/paginator"

type GetUserDetailRequest struct {
	UserID int64 `params:"id" validate:"required"`
}

type GetUserDetailReply struct {
	UserDetail
}

type GetAllUsersRequest struct {
	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
	// query condition
	QueryCond string `query:"cond" validate:"omitempty,oneof=newest" label:"排序"`
}

type GetAllUsersReply struct {
	paginator.Page[UserDetail]
}

type UserLoginRequset struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginReply struct {
	User        BasicUserInfo
	AccessToken string
}

type UserRegisterRequest struct {
	Name     string `json:"name" validate:"required,gte=3,lte=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterReply struct {
	User        BasicUserInfo
	AccessToken string
}
