package services

import (
	"fiberLearn/model"
	"fiberLearn/pkg/errcode"
	"fiberLearn/pkg/snowflake"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegistService struct {
	Name     string `json:"name" validate:"required,gte=3,lte=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginService struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func GetUser(userID int) (*model.UserDetail, error) {
	var user model.UserDetail
	if err := model.DB.Table("users").Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUsers(offset, limit int) ([]*model.UserInfo, int64, error) {
	users := []*model.UserInfo{}
	if err := model.DB.Table("users").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	var total int64
	if err := model.DB.Table("users").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// ValidUsername 验证用户
func ValidUsername(username string) error {
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return errcode.UsernameCharLimit
	}
	return nil
}

func (urs *UserRegistService) Regist() (*model.UserDetail, error) {
	if exist, err := model.IsUserExist(urs.Email); err != nil {
		return nil, err
	} else if exist {
		return nil, nil
	}
	if err := ValidUsername(urs.Name); err != nil {
		return nil, err
	}
	user := model.User{
		UserID: snowflake.GenID(),
		Name:   urs.Name,
		Email:  urs.Email,
	}
	if p, err := HashAndSalt(urs.Password); err != nil {
		return nil, err
	} else {
		user.Password = p
	}
	if err := model.DB.Table("users").Create(&user).Error; err != nil {
		return nil, err
	}
	return &model.UserDetail{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func HashAndSalt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func CheckPassword(userPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
	return err == nil
}

func (uls *UserLoginService) Login() (*model.UserDetail, string, error) {
	var user model.User
	if err := model.DB.Table("users").Where("email = ?", uls.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errcode.UnauthorizedAuthNotExist
		}
		return nil, "", errcode.UserLoginFailed
	}
	if !CheckPassword(user.Password, uls.Password) {
		return nil, "", errcode.UnauthorizedAuthFailed
	}
	// Create the Claims
	claims := jwt.MapClaims{
		"name": user.Name,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, "", errcode.UnauthorizedTokenGenerate
	}
	return &model.UserDetail{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, t, nil
}
