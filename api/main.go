package main

import (
	"awesomeProject/api/internal/endpoints"
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/repositories"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Testing Portfolio",
		AppName:      "Portfolio v.0.0.1",
		ErrorHandler: endpoints.ErrorHandler,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test index")
	})
	endpoints.SetupRoutes(app)
	err := app.Listen(":8080")

	if err != nil {
		log.Fatalln("Can't listen to port 8080 or app can't start.")
	}

	conn, err := repositories.CreateConnection()
	err = repositories.FirstInitialization(conn, models.Schema)
}
