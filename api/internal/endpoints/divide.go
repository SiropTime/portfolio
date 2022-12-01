package endpoints

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"go/types"
	"log"
	"math/big"
	"portfolioTask/api/internal/crud"
	"portfolioTask/api/internal/etc"
	"portfolioTask/api/internal/models"
	"strconv"
)

func getTokenPrice(token models.TokenInPortfolio) (string, error) {
	tokenAmount, _, err := big.ParseFloat(token.Amount, 10, token.Decimals, big.AwayFromZero)
	if err != nil {
		return "", err
	}
	log.Println(tokenAmount.String())
	return "", nil
}

func sumTokensPriceInPortfolio(portfolio models.PortfolioResponse, tokensPrices map[string]string) (*models.TokensInfo, error) {
	nativeToken, err := etc.GetNativeTokenInfo(portfolio.ChainId)
	if err != nil {
		return nil, err
	}
	tokensInfo := models.TokensInfo{NativeToken: nativeToken}
	sum := decimal.NewFromInt(0)
	for _, token := range portfolio.Tokens {
		_price, success := big.NewInt(0).SetString(tokensPrices[token.Address], 10)
		if !success {
			return nil, types.Error{Msg: "Not valid price type"}
		}
		realPrice := decimal.NewFromBigInt(_price, -int32(nativeToken.TokenDecimals))
		_amount, success := big.NewInt(0).SetString(token.Amount, 10)
		realAmount := decimal.NewFromBigInt(_amount, -int32(token.Decimals))
		_temp := realPrice.Mul(realAmount)
		tokensInfo.Tokens = append(tokensInfo.Tokens, models.TokenRealData{
			Ticker:     token.Ticker,
			Address:    token.Address,
			TotalPrice: _temp,
		})
		sum = sum.Add(_temp)
	}
	tokensInfo.Total = sum
	return &tokensInfo, nil
}

func calculatePortfolioProportions(portfolio models.PortfolioResponse) (*models.PortfolioProportionsResponse, error) {
	tokensPrices, err := etc.GetTokensPrices(portfolio.ChainId)
	portfolioResponse := models.PortfolioProportionsResponse{
		Name:    portfolio.Name,
		Id:      portfolio.Id,
		ChainId: portfolio.ChainId,
	}
	if err != nil {
		return nil, err
	}
	tokensInfo, err := sumTokensPriceInPortfolio(portfolio, tokensPrices)
	if err != nil {
		return nil, err
	}

	for _, token := range tokensInfo.Tokens {
		proportion := token.TotalPrice.Div(tokensInfo.Total)

		portfolioResponse.TokensProportions = append(portfolioResponse.TokensProportions, models.TokenProportions{
			Ticker:     token.Ticker,
			Address:    token.Address,
			Proportion: proportion,
		})
	}

	return &portfolioResponse, nil

}

func calculatePortfolioWithAmount(portfolio models.PortfolioProportionsResponse,
	amount string, fromAddress string, gasPrice int) (*models.PortfolioAfterQuote, error) {
	_amount, success := new(big.Int).SetString(amount, 10)
	if !success {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Amount is represented in not valid type. Try another value.",
		}
	}
	//tokensPrices, err := etc.GetTokensPrices(portfolio.ChainId)
	//if err != nil {
	//	return err
	//}
	tokenFromAPI, err := crud.GetTokenDetails(portfolio.ChainId, fromAddress)
	if err != nil {
		return nil, err
	}
	//_price, success := new(big.Int).SetString(tokensPrices[fromAddress], 10)
	//if !success {
	//	return err
	//}
	nativeToken, err := etc.GetNativeTokenInfo(portfolio.ChainId)

	fromRealAmount := decimal.NewFromBigInt(_amount, -int32(tokenFromAPI.TokenDecimals))
	portfolioResponse := models.PortfolioAfterQuote{
		Id:      portfolio.Id,
		ChainId: portfolio.ChainId,
		Name:    portfolio.Name,
	}
	for _, token := range portfolio.TokensProportions {
		amountForCurrentToken := fromRealAmount.Mul(token.Proportion)
		queryForQuote := models.QuoteQuery{
			ChainId:          portfolio.ChainId,
			GasPrice:         gasPrice,
			Amount:           amountForCurrentToken.Shift(int32(nativeToken.TokenDecimals)).BigInt().String(),
			FromTokenAddress: fromAddress,
			ToTokenAddress:   token.Address,
		}
		if !(queryForQuote.FromTokenAddress == queryForQuote.ToTokenAddress) {
			quoteResult, e := etc.GetQuoteApi(queryForQuote)
			if e != nil {
				return nil, e
			}
			portfolioResponse.Tokens = append(portfolioResponse.Tokens, models.TokenQuote{
				FinalAmount:  quoteResult.ToTokenAmount,
				EstimatedGas: quoteResult.EstimatedGas,
				Address:      token.Address,
				Ticker:       token.Ticker,
			})
		} else {
			portfolioResponse.Tokens = append(portfolioResponse.Tokens, models.TokenQuote{
				FinalAmount:  queryForQuote.Amount,
				EstimatedGas: 0,
				Address:      token.Address,
				Ticker:       token.Ticker,
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
	portfolio, err := crud.ReadPortfolio(pId)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return err
	}

	portfolioResponse, err := calculatePortfolioProportions(portfolio)
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
	portfolio, err := crud.ReadPortfolio(pId)
	if err != nil {
		return &fiber.Error{Code: 404, Message: "This portfolio doesn't exist"}
	}
	calculatedPortfolio, err := calculatePortfolioProportions(portfolio)
	if err != nil {
		return err
	}

	quotePortfolio, err := calculatePortfolioWithAmount(
		*calculatedPortfolio, amount,
		contractAddress, gasPrice,
	)
	if err != nil {
		return err
	}

	return c.JSON(quotePortfolio)
}
