package endpoints

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	app.Get("/portfolios/:id", GetPortfolio)
	app.Get("/portfolios", GetAllPortfolios)
}
