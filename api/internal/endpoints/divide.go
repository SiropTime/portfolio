package endpoints

import (
	"awesomeProject/api/internal/crud"
	"awesomeProject/api/internal/etc"
	"awesomeProject/api/internal/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"go/types"
	"log"
	"math/big"
	"strconv"
)

func calculateRealAmount(token models.TokenInPortfolio) {

}

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
	//c.SendString(sum.String())
	//return nil
	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portfolioResponse)
}
