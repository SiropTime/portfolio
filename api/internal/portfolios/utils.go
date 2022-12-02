package portfolios

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/shopspring/decimal"
	"go/types"
	"math/big"
	tokens "portfolioTask/api/internal/token"
	"portfolioTask/api/internal/token/repository"
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

func CalculatePortfolioWithAmount(portfolio ProportionsResponsePortfolio,
	amount string, fromAddress string, gasPrice int) (*AfterQuotePortfolio, error) {
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
	portfolioResponse := AfterQuotePortfolio{
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
				NativePrice:  queryForQuote.Amount,
				EstimatedGas: quoteResult.EstimatedGas,
				Address:      proportionToken.Address,
				Ticker:       proportionToken.Ticker,
			})
		} else {
			portfolioResponse.Tokens = append(portfolioResponse.Tokens, tokens.QuoteToken{
				FinalAmount:  queryForQuote.Amount,
				NativePrice:  queryForQuote.Amount,
				EstimatedGas: 0,
				Address:      proportionToken.Address,
				Ticker:       proportionToken.Ticker,
			})
		}

	}
	return &portfolioResponse, nil
}
