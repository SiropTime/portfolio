package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"log"
	"os"
	errors "portfolioTask/api/internal/httpServer/error"
	"portfolioTask/api/internal/httpServer/utils"
	"portfolioTask/api/internal/portfolios/delivery"
	_ "portfolioTask/docs"
)

// @title Portfolio API
// @version 1.0
// @description This is a portfolio API server. There are endpoints for creating, updating, deleting, getting portfolios and making swaps and quotes with them.
// @contact.name API developer
// @contact.url https://t.me/KlenoviySIr
// @contact.email KlenoviySir@yandex.ru

// @host localhost:8080
// @BasePath /portfolios
// @schemes http
func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Testing Portfolio",
		AppName:      "Portfolio v.1.0.0",
		ErrorHandler: errors.ErrorHandler,
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Get("/", utils.HealthCheck)
	app.Get("/docs/*", swagger.HandlerDefault)
	delivery.SetupRoutes(app)
	err := app.Listen(os.Getenv("APP_PORT"))

	if err != nil {
		log.Fatalln("Can't listen to port 8080 or app can't start.")
	}
}
