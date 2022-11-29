package main

import (
	"awesomeProject/api/internal/endpoints"
	"awesomeProject/api/pkg/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"log"
)

var DBConn *sqlx.DB

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Testing Portfolio",
		AppName:      "Portfolio v.0.0.1",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test index")
	})
	endpoints.SetupRoutes(app)
	err := app.Listen(":8080")

	if err != nil {
		log.Fatalln("Can't listen to port 8080 or app can't start.")
	}

	DBConn, err = repositories.CreateConnection()

}
