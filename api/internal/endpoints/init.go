package endpoints

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/portfolios/:id", GetPortfolio)
	app.Get("/portfolios", GetAllPortfolios)
	app.Post("/portfolios", PostPortfolio)
	app.Delete("/portfolios/:id", DeletePortfolio)
	app.Patch("/portfolios/:id/addTokens", AddNewTokensToPortfolio)
	app.Put("/portfolios/:id", UpdatePortfolio)
	app.Get("/portfolios/:id/proportions", GetPortfolioProportions)
	app.Get("/portfolios/:id/count", GetCountedPortfolio)
}

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	// Default status cde
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return ctx.Status(code).SendString(err.Error())
}
