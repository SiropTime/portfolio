package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	errors "portfolioTask/api/internal/httpServer/error"
	"portfolioTask/api/internal/portfolios/delivery"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Testing Portfolio",
		AppName:      "Portfolio v.1.0.0",
		ErrorHandler: errors.ErrorHandler,
	})
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test index")
	})
	delivery.SetupRoutes(app)
	err := app.Listen(os.Getenv("APP_PORT"))

	if err != nil {
		log.Fatalln("Can't listen to port 8080 or app can't start.")
	}
}
