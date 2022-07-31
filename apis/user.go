package apis

import (
	"fiberLearn/model"
	"fiberLearn/pkg/snowflake"
	"fiberLearn/services"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func GetAllUsers(c *fiber.Ctx) error {
	var users []services.User
	if err := model.DB.Find(&users).Error; err != nil {
		return err
	}
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user model.User
	if err := model.DB.First(&user, id).Error; err != nil {
		return err
	}
	return c.JSON(user)
}

func Regist(c *fiber.Ctx) error {
	var service services.UserRegistService

	if err := c.BodyParser(&service); err != nil {
		return err
	}
	user := model.User{
		UserID: snowflake.GenID(),
		Name:   service.Name,
		Email:  service.Email,
	}
	if err := user.HashAndSalt(service.Password); err != nil {
		return err
	}
	if exist, err := user.IsExist(); err != nil {
		return err
	} else if !exist {
		if err := model.DB.Create(&user).Error; err != nil {
			fmt.Println(err == gorm.ErrRegistered)
			return err
		}
	} else {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"Error": "User already exists"})
	}
	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var service services.UserLoginService
	if err := c.BodyParser(&service); err != nil {
		return err
	}
	var user model.User
	if err := model.DB.Where("email = ?", service.Email).First(&user).Error; err != nil {
		return err
	}
	if !user.CheckPassword(service.Password) {
		return fiber.NewError(fiber.StatusUnauthorized, "Wrong password")
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
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t, "user": user})
}
