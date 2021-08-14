package main

import (
	"app/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	routes.Serve(app)
	
	app.Listen(":" + os.Getenv("PORT"))
}
