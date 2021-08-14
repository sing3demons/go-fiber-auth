package routes

import (
	"app/controllers"
	"app/database"

	"github.com/gofiber/fiber/v2"
)

func Serve(app *fiber.App) {
	db := database.GetDB()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	v1 := app.Group("/api/v1")
	authController := controllers.Auth{DB: db}

	{
		v1.Post("/register", authController.Register)
		v1.Post("/login", authController.Login)
		v1.Get("/user", authController.User)
		v1.Get("/logout", authController.Logout)
	}

	userController := controllers.ResetPassword{DB: db}

	{
		v1.Post("/forgot", userController.Forgot)
		v1.Post("/reset", userController.Reset)
	}

}
