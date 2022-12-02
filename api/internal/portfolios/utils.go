package portfolios

import (
	"github.com/shopspring/decimal"
	"go/types"
	"math/big"
	tokens "portfolioTask/api/internal/token"
	"portfolioTask/api/pkg/clients/swapAPI"
)

func SumTokensPriceInPortfolio(portfolio ResponsePortfolio, tokensPrices map[string]string) (*swapAPI.TokensInfo, error) {
	nativeToken, err := swapAPI.GetNativeTokenInfo(portfolio.ChainId)
	if err != nil {
		return nil, err
	}
	tokensInfo := swapAPI.TokensInfo{NativeToken: nativeToken}
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
		tokensInfo.Tokens = append(tokensInfo.Tokens, tokens.RealDataToken{
			Ticker:     token.Ticker,
			Address:    token.Address,
			TotalPrice: _temp,
		})
		sum = sum.Add(_temp)
	}
	tokensInfo.Total = sum
	return &tokensInfo, nil
}

func CalculatePortfolioProportions(portfolio ResponsePortfolio) (*ProportionsResponsePortfolio, error) {
	tokensPrices, err := swapAPI.GetTokensPrices(portfolio.ChainId)
	portfolioResponse := ProportionsResponsePortfolio{
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

		portfolioResponse.TokensProportions = append(portfolioResponse.TokensProportions, tokens.ProportionsToken{
			Ticker:     token.Ticker,
			Address:    token.Address,
			Proportion: proportion,
		})
	}

	return &portfolioResponse, nil

}
