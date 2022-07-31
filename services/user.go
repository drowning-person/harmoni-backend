package services

import "time"

type User struct {
	UserID    int64      `json:"user_id,string"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_time"`
}

type UserRegistService struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginService struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
