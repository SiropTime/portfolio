package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {

	app := fiber.New(fiber.Config{
		ServerHeader: "Testing Portfolio",
		AppName:      "Portfolio v.0.0.1",
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("test index")
	})

	err := app.Listen(":8080")

	if err != nil {
		log.Fatal("Can't listen to port 8080 or app can't start.")
	}

}
