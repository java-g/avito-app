package main

import (
	"avito-app/database"
	"avito-app/middlewares"
	"avito-app/models"
	"avito-app/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var config models.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config = models.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		TimeZone: os.Getenv("DB_TIMEZONE"),
		JwtKey:   os.Getenv("JWTKEY"),
	}
}

func main() {
	database.InitDB(config)

	database.Migration()
	database.MigrationData()

	middlewares.SetJWT(config.JwtKey)
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05.000000",
		TimeZone:   "Europe/Moscow",
	}))
	routes.Setup(app)
	app.Listen(":8000")
}
