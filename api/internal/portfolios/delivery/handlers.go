package delivery

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"portfolioTask/api/internal/portfolios"
	"portfolioTask/api/internal/portfolios/repository"
	tokens "portfolioTask/api/internal/token"
	"strconv"
)

func GetPortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	portRes, err := token.ReadPortfolio(pId)
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portRes)
}

func GetAllPortfolios(c *fiber.Ctx) error {
	portfoliosResponse, err := token.ReadAllPortfolios()
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portfoliosResponse)
}

func PostPortfolio(c *fiber.Ctx) error {
	portfolio := new(portfolios.InputPortfolio)
	if err := c.BodyParser(portfolio); err != nil {
		c.Status(503)
		return err
	}
	portfolioResponse, err := token.CreatePortfolio(*portfolio)
	if err != nil {
		return err
	}
	c.Status(201)
	err = c.JSON(portfolioResponse)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't marshal portfolios to JSON",
		}
	}
	return nil
}

func DeletePortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = token.DeletePortfolio(pId)
	if err != nil {
		return err
	}
	return nil
}

func AddNewTokensToPortfolio(c *fiber.Ctx) error {
	tokensList := new(tokens.InputTokens)
	if err := c.BodyParser(tokensList); err != nil {
		return fiber.ErrBadRequest
	}
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	err = token.AddNewTokens(pId, *tokensList)
	if err != nil {
		return err
	}

	_p, err := token.ReadPortfolio(pId)
	if err != nil {
		return &fiber.Error{
			Code:    404,
			Message: "Can't find portfolios with this id",
		}
	}
	portfolioResponse, err := portfolios.CalculatePortfolioProportions(_p)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't calculate new proportions for portfolios",
		}
	}

	c.Status(fiber.StatusAccepted)
	err = c.JSON(portfolioResponse)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't marshal portfolios to JSON",
		}
	}

	return nil
}

func UpdatePortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	portfolio := new(portfolios.InputPortfolio)
	if err := c.BodyParser(portfolio); err != nil {
		return fiber.ErrBadRequest
	}
	if portfolio.ChainId == 0 || portfolio.Tokens == nil || len(portfolio.Name) == 0 {
		return fiber.ErrBadRequest
	}
	err = token.UpdatePortfolio(pId, *portfolio)
	if err != nil {
		return err
	}

	return nil
}

func GetPortfolioProportions(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}
	portfolio, err := token.ReadPortfolio(pId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return err
	}

	portfolioResponse, err := portfolios.CalculatePortfolioProportions(portfolio)
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portfolioResponse)
}

func GetCountedPortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return err
	}
	amount := c.Query("amount")
	contractAddress := c.Query("contractAddress")
	gasPrice, err := strconv.Atoi(c.Query("gasPrice"))
	if len(amount) == 0 || len(contractAddress) == 0 {
		return &fiber.Error{Code: 404, Message: "Query parameters are not found. Check documentation"}
	}
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "GasPrice is not represented correctly. Try another value",
		}
	}
	portfolio, err := token.ReadPortfolio(pId)
	if err != nil {
		return &fiber.Error{Code: 404, Message: "This portfolios doesn't exist"}
	}
	calculatedPortfolio, err := portfolios.CalculatePortfolioProportions(portfolio)
	if err != nil {
		return err
	}

	quotePortfolio, err := portfolios.CalculatePortfolioWithAmount(
		*calculatedPortfolio, amount,
		contractAddress, gasPrice,
	)
	if err != nil {
		return err
	}

	return c.JSON(quotePortfolio)
}
