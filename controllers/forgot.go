package controllers

import (
	"app/models"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Reset struct {
	DB *gorm.DB
}

func (r *Reset) Forgot(ctx *fiber.Ctx) error {
	var form map[string]string
	if err := ctx.BodyParser(&form); err != nil {
		return err
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{Email: form["email"],
		Token: token,
	}

	r.DB.Create(&passwordReset)

	return ctx.JSON(passwordReset)
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
