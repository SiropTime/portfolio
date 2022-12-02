package delivery

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"math/big"
	"portfolioTask/api/internal/portfolios"
	"portfolioTask/api/internal/portfolios/repository"
	tokens "portfolioTask/api/internal/token"
	"portfolioTask/api/internal/token/repository"
	"portfolioTask/api/pkg/clients/swapAPI"
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

func CalculatePortfolioWithAmount(portfolio portfolios.ProportionsResponsePortfolio,
	amount string, fromAddress string, gasPrice int) (*portfolios.AfterQuotePortfolio, error) {
	_amount, success := new(big.Int).SetString(amount, 10)
	if !success {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Amount is represented in not valid type. Try another value.",
		}
	}

	tokenFromAPI, err := repository.GetTokenDetails(portfolio.ChainId, fromAddress)
	if err != nil {
		return nil, err
	}

	nativeToken, err := swapAPI.GetNativeTokenInfo(portfolio.ChainId)

	fromRealAmount := decimal.NewFromBigInt(_amount, -int32(tokenFromAPI.TokenDecimals))
	portfolioResponse := portfolios.AfterQuotePortfolio{
		Id:      portfolio.Id,
		ChainId: portfolio.ChainId,
		Name:    portfolio.Name,
	}
	for _, proportionToken := range portfolio.TokensProportions {
		amountForCurrentToken := fromRealAmount.Mul(proportionToken.Proportion)
		queryForQuote := swapAPI.QuoteQuery{
			ChainId:          portfolio.ChainId,
			GasPrice:         gasPrice,
			Amount:           amountForCurrentToken.Shift(int32(nativeToken.TokenDecimals)).BigInt().String(),
			FromTokenAddress: fromAddress,
			ToTokenAddress:   proportionToken.Address,
		}
		fmt.Printf("Address: %s; AmountForQuote: %s; RealAmount: %s; Ticker: %s\n", queryForQuote.ToTokenAddress, queryForQuote.Amount, amountForCurrentToken.String(), proportionToken.Ticker)
		if !(queryForQuote.FromTokenAddress == queryForQuote.ToTokenAddress) {
			quoteResult, e := swapAPI.GetQuoteApi(queryForQuote)
			if e != nil {
				return nil, e
			}
			portfolioResponse.Tokens = append(portfolioResponse.Tokens, tokens.QuoteToken{
				FinalAmount:  quoteResult.ToTokenAmount,
				EstimatedGas: quoteResult.EstimatedGas,
				Address:      proportionToken.Address,
				Ticker:       proportionToken.Ticker,
			})
		} else {
			portfolioResponse.Tokens = append(portfolioResponse.Tokens, tokens.QuoteToken{
				FinalAmount:  queryForQuote.Amount,
				EstimatedGas: 0,
				Address:      proportionToken.Address,
				Ticker:       proportionToken.Ticker,
			})
		}

	}
	return &portfolioResponse, nil
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

	quotePortfolio, err := CalculatePortfolioWithAmount(
		*calculatedPortfolio, amount,
		contractAddress, gasPrice,
	)
	if err != nil {
		return err
	}

	return c.JSON(quotePortfolio)
}
