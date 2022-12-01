package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"portfolioTask/api/internal/endpoints"
	"portfolioTask/api/internal/models"
	"portfolioTask/api/pkg/repositories"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Testing Portfolio",
		AppName:      "Portfolio v.0.0.2",
		ErrorHandler: endpoints.ErrorHandler,
	})
	app.Use(logger.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test index")
	})
	endpoints.SetupRoutes(app)
	err := app.Listen(":8080")

	if err != nil {
		log.Fatalln("Can't listen to port 8080 or app can't start.")
	}

	err = repositories.FirstInitialization(models.Schema)
}
