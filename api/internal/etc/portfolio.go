package etc

import (
	"github.com/shopspring/decimal"
	"go/types"
	"math/big"
	"portfolioTask/api/internal/models"
)

func SumTokensPriceInPortfolio(portfolio models.PortfolioResponse, tokensPrices map[string]string) (*models.TokensInfo, error) {
	nativeToken, err := GetNativeTokenInfo(portfolio.ChainId)
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

func CalculatePortfolioProportions(portfolio models.PortfolioResponse) (*models.PortfolioProportionsResponse, error) {
	tokensPrices, err := GetTokensPrices(portfolio.ChainId)
	portfolioResponse := models.PortfolioProportionsResponse{
		Name:    portfolio.Name,
		Id:      portfolio.Id,
		ChainId: portfolio.ChainId,
	}
	if err != nil {
		return nil, err
	}
	tokensInfo, err := SumTokensPriceInPortfolio(portfolio, tokensPrices)
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
