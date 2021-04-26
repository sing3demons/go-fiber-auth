package controllers

import (
	"app/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type CreateUser struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
type loginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	jwt.StandardClaims
}

func (a *Auth) Register(ctx *fiber.Ctx) error {
	var user models.User
	var form CreateUser
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if form.Password != form.PasswordConfirm {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Password do not macth"})
	}

	copier.Copy(&user, &form)
	user.Password = user.GenerateFromPassword()
	a.DB.Create(&user)

	var serializedUser UserResponse
	copier.Copy(&serializedUser, &user)
	return ctx.Status(fiber.StatusCreated).JSON(serializedUser)
}

func (a *Auth) Login(ctx *fiber.Ctx) error {
	var form loginUser
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var user models.User

	if err := a.DB.Where("email = ?", form.Email).First(&user).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})

	}

	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Local().Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte("32d4b883-fcdd-4d63-8091-1a232218448a"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{"token": token})
}

func (a *Auth) User(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("token")

	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte("32d4b883-fcdd-4d63-8091-1a232218448a"), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}
	claims := token.Claims.(*Claims)
	id := claims.Issuer

	var user models.User
	a.DB.Where("id = ?", id).First(&user)
	var serializedUser UserResponse
	copier.Copy(&serializedUser, &user)
	return ctx.JSON(fiber.Map{"user": serializedUser})
}

func (a *Auth) Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	ctx.Cookie(&cookie)
	return ctx.JSON(fiber.Map{"message": "success"})
}


