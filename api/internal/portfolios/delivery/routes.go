package delivery

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {

	app.Get("/portfolios/:id", GetPortfolio)
	app.Get("/portfolios", GetAllPortfolios)
	app.Post("/portfolios", PostPortfolio)
	app.Delete("/portfolios/:id", DeletePortfolio)
	app.Patch("/portfolios/:id/addTokens", AddNewTokensToPortfolio)
	app.Put("/portfolios/:id", UpdatePortfolio)
	app.Get("/portfolios/:id/proportions", GetPortfolioProportions)
	app.Get("/portfolios/:id/count", GetCountedPortfolio)
	//app.Get("/portfolios/:id/transactions", GetTransactions)
}
