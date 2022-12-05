package delivery

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
	portfoliosAPI := app.Group("/portfolios")
	portfoliosAPI.Get("/:id", GetPortfolio)
	portfoliosAPI.Get("/", GetAllPortfolios)
	portfoliosAPI.Post("/", PostPortfolio)
	portfoliosAPI.Delete("/:id", DeletePortfolio)
	portfoliosAPI.Patch("/:id/addTokens", AddNewTokensToPortfolio)
	portfoliosAPI.Put("/:id", UpdatePortfolio)
	portfoliosAPI.Get("/:id/proportions", GetPortfolioProportions)
	portfoliosAPI.Get("/:id/count", GetCountedPortfolio)
	portfoliosAPI.Get("/:id/transactions", GetTransactions)
}
