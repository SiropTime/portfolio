package endpoints

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"portfolioTask/api/internal/crud"
	"portfolioTask/api/internal/etc"
	"portfolioTask/api/internal/models"
	"strconv"
)

func GetPortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	portRes, err := crud.ReadPortfolio(pId)
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portRes)
}

func GetAllPortfolios(c *fiber.Ctx) error {
	portfolios, err := crud.ReadAllPortfolios()
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portfolios)
}

func PostPortfolio(c *fiber.Ctx) error {
	portfolio := new(models.PortfolioInput)
	if err := c.BodyParser(portfolio); err != nil {
		c.Status(503)
		return err
	}
	portfolioResponse, err := crud.CreatePortfolio(*portfolio)
	if err != nil {
		return err
	}
	c.Status(201)
	err = c.JSON(portfolioResponse)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't marshal portfolio to JSON",
		}
	}
	return nil
}

func DeletePortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = crud.DeletePortfolio(pId)
	if err != nil {
		return err
	}
	return nil
}

func AddNewTokensToPortfolio(c *fiber.Ctx) error {
	tokens := new(models.TokensInput)
	if err := c.BodyParser(tokens); err != nil {
		return fiber.ErrBadRequest
	}
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	err = crud.AddNewTokens(pId, *tokens)
	if err != nil {
		return err
	}

	_p, err := crud.ReadPortfolio(pId)
	if err != nil {
		return &fiber.Error{
			Code:    404,
			Message: "Can't find portfolio with this id",
		}
	}
	portfolioResponse, err := etc.CalculatePortfolioProportions(_p)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't calculate new proportions for portfolio",
		}
	}

	c.Status(fiber.StatusAccepted)
	err = c.JSON(portfolioResponse)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't marshal portfolio to JSON",
		}
	}

	return nil
}

func UpdatePortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	portfolio := new(models.PortfolioInput)
	if err := c.BodyParser(portfolio); err != nil {
		return fiber.ErrBadRequest
	}
	if portfolio.ChainId == 0 || portfolio.Tokens == nil || len(portfolio.Name) == 0 {
		return fiber.ErrBadRequest
	}
	err = crud.UpdatePortfolio(pId, *portfolio)
	if err != nil {
		return err
	}

	return nil
}
