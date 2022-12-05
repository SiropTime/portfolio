package delivery

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"portfolioTask/api/internal/portfolios"
	rep "portfolioTask/api/internal/portfolios/repository"
	tokens "portfolioTask/api/internal/token"
	"portfolioTask/api/pkg/clients/swapAPI"
	"strconv"
)

// GetPortfolio godoc
// @Summary Get portfolio by id
// @Description Get portfolio by id with amount of tokens
// @Tags portfolio
// @Accept  */*
// @Param id path int true "Portfolio id"
// @Produce  json
// @Success 200 {object} portfolios.ResponsePortfolio
// @Failure 404
// @Failure 500
// @Router /{id} [get]
func GetPortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	portRes, err := rep.ReadPortfolio(pId)
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portRes)
}

// GetAllPortfolios godoc
// @Summary Get all portfolios in DB
// @Description Get all portfolios in DB with their tokens with amount
// @Tags portfolio
// @Accept  */*
// @Produce  json
// @Success 200 {object} []portfolios.ResponsePortfolio
// @Failure 404
// @Failure 500
// @Router / [get]
func GetAllPortfolios(c *fiber.Ctx) error {
	portfoliosResponse, err := rep.ReadAllPortfolios()
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portfoliosResponse)
}

// PostPortfolio godoc
// @Summary Create new portfolio
// @Description Create new portfolio with tokens
// @Tags portfolio
// @Accept  application/json
// @Param portfolio body portfolios.InputPortfolio true "Portfolio"
// @Produce  json
// @Success 201 {object} portfolios.ProportionsResponsePortfolio
// @Failure 400
// @Failure 500
// @Router / [post]
func PostPortfolio(c *fiber.Ctx) error {
	portfolio := new(portfolios.InputPortfolio)
	if err := c.BodyParser(portfolio); err != nil {
		c.Status(503)
		return err
	}
	portfolioResponse, err := rep.CreatePortfolio(*portfolio)
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

// DeletePortfolio godoc
// @Summary Delete portfolio by id
// @Description Delete portfolio with its tokens by id
// @Tags portfolio
// @Accept  */*
// @Param id path int true "Portfolio id"
// @Produce  json
// @Success 204
// @Failure 404
// @Failure 500
// @Router /{id} [delete]
func DeletePortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = rep.DeletePortfolio(pId)
	if err != nil {
		return err
	}
	return nil
}

// AddNewTokensToPortfolio godoc
// @Summary Add new tokens to portfolio
// @Description Add new tokens to portfolio by id
// @Tags portfolio
// @Accept  application/json
// @Param id path int true "Portfolio id"
// @Param tokens body token.InputTokens true "Tokens"
// @Produce  json
// @Success 202 {object} portfolios.ProportionsResponsePortfolio
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id}/tokens [patch]
func AddNewTokensToPortfolio(c *fiber.Ctx) error {
	tokensList := new(tokens.InputTokens)
	if err := c.BodyParser(tokensList); err != nil {
		return fiber.ErrBadRequest
	}
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.ErrBadRequest
	}
	err = rep.AddNewTokens(pId, *tokensList)
	if err != nil {
		return err
	}

	_p, err := rep.ReadPortfolio(pId)
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

// UpdatePortfolio godoc
// @Summary Update portfolio
// @Description Update whole portfolio by id
// @Tags portfolio
// @Accept  application/json
// @Param id path int true "Portfolio id"
// @Param portfolio body portfolios.InputPortfolio true "Portfolio"
// @Produce  json
// @Success 202 {object} portfolios.ProportionsResponsePortfolio
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id} [put]
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
	err = rep.UpdatePortfolio(pId, *portfolio)
	if err != nil {
		return err
	}

	return nil
}

// GetPortfolioProportions godoc
// @Summary Get portfolio proportions
// @Description Get portfolio with tokens represented with proportions by id
// @Tags portfolio
// @Accept  */*
// @Param id path int true "Portfolio id"
// @Produce  json
// @Success 200 {object} portfolios.ProportionsResponsePortfolio
// @Failure 404
// @Failure 500
// @Router /{id}/proportions [get]
func GetPortfolioProportions(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}
	portfolio, err := rep.ReadPortfolio(pId)
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

// GetCountedPortfolio godoc
// @Summary Get portfolio after quote
// @Description Get portfolio with pre-calculated tokens values for transaction and gas by portfolio id and token with amount from which will transaction be made
// @Tags portfolio
// @Accept  */*
// @Param id path int true "Portfolio id"
// @Param amount query string true "Amount of token"
// @Param contractAddress query string true "Contract address of token"
// @Param gasPrice query string true "Gas price for transaction to count estimated gas of this chain"
// @Produce  json
// @Success 200 {object} portfolios.AfterQuotePortfolio
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id}/count [get]
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
	portfolio, err := rep.ReadPortfolio(pId)
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

// GetTransactions godoc
// @Summary Get transactions
// @Description Get transactions data for metamask for current portfolio with amount of token and wallet address
// @Tags portfolio
// @Accept  */*
// @Param id path int true "Portfolio id"
// @Param amount query string true "Amount of token"
// @Param fromTokenAddress query string true "Contract address of token"
// @Param fromAddress query string true "User's wallet address"
// @Param slippage query string true "Slippage for transaction"
// @Produce  json
// @Success 200 {object} portfolios.AfterSwapPortfolio
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{id}/transactions [get]
func GetTransactions(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	amount := c.Query("amount")
	fromAddress := c.Query("fromAddress")
	fromTokenAddress := c.Query("fromTokenAddress")
	_slippage := c.Query("slippage")
	slippage, err := strconv.Atoi(_slippage)
	if len(amount) == 0 || len(fromAddress) == 0 {
		return &fiber.Error{Code: 404, Message: "Query parameters are not found. Check documentation"}
	}
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Slippage is not represented correctly. Try another value",
		}
	}
	portfolioToParse, err := rep.ReadPortfolio(pId)
	if err != nil {
		return err
	}
	portfolioToSend, err := portfolios.CalculatePortfolioProportions(portfolioToParse)
	if err != nil {
		return err
	}
	portfolioResponse, err := portfolios.FormTransaction(portfolioToSend,
		swapAPI.SwapQuery{
			FromTokenAddress: fromTokenAddress,
			FromAddress:      fromAddress,
			Amount:           amount,
			Slippage:         slippage,
			ChainId:          portfolioToSend.ChainId,
			ToTokenAddress:   "",
		})
	if err != nil {
		return err
	}

	return c.JSON(portfolioResponse)
}
