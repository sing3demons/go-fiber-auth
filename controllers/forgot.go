package controllers

import (
	"app/models"
	"math/rand"

	"net/smtp"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ResetPassword struct {
	DB *gorm.DB
}

func (r *ResetPassword) Forgot(ctx *fiber.Ctx) error {
	var form map[string]string
	if err := ctx.BodyParser(&form); err != nil {
		return err
	}

	token := RandStringRunes(12)
	passwordReset := models.PasswordReset{Email: form["email"],
		Token: token,
	}

	r.DB.Create(&passwordReset)

	from := "sing@dev.com"
	to := []string{form["email"]}

	url := "http://localhost:8080/api/v1/reset/" + token

	message := []byte("Click <a href=\"" + url + "\">here</a> to reset your password!")
	err := smtp.SendMail("0.0.0.0:1025", nil, from, to, message)
	if err != nil {
		return err
	}

	return ctx.JSON(passwordReset)
}

func (r *ResetPassword) Reset(ctx *fiber.Ctx) error {
	var form map[string]string
	if err := ctx.BodyParser(&form); err != nil {
		return err
	}

	if form["password"] != form["password_confirm"] {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Password do not macth"})
	}

	passwordReset := models.PasswordReset{}

	if err := r.DB.Where("token = ?", form["token"]).Last(&passwordReset).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid token!"})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(form["password"]), 14)
	r.DB.Model(&models.User{}).Where("email = ?", passwordReset.Email).Update("password", password)

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}

func RandStringRunes(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
