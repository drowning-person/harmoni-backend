package services

type UserRegistService struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginService struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
