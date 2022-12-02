package swapAPI

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"portfolioTask/api/internal/cconst"
)

func GetNativeTokenInfo(chainId int) (TokenAPI, error) {
	tokens, err := GetTokensAPI(chainId)
	if err != nil {
		return TokenAPI{}, err
	}
	for _, token := range tokens {
		if token.TokenContractAddress == cconst.NativeAddress {
			return token, nil
		}
	}
	return TokenAPI{}, fiber.ErrNotFound
}

func GetTokensAPI(chainId int) ([]TokenAPI, error) {
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(cconst.SwapAPIURL+"/tokens?chainId=%d", chainId))
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}
	var tokensBody TokenRequestAPI
	err = json.Unmarshal(res.Body(), &tokensBody)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't unmarshal body from external API in needed structure",
		}
	}

	if !tokensBody.Success {
		return nil, &fiber.Error{
			Code:    tokensBody.StatusCode,
			Message: "Got an error from external API. See code",
		}
	}

	return tokensBody.Result, nil
}

func GetTokensPrices(chainId int) (map[string]string, error) {
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(cconst.SwapAPIURL+"/prices?chainId=%d", chainId))
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}
	var tokensBody TokenPriceAPI

	err = json.Unmarshal(res.Body(), &tokensBody)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't unmarshal body from external API in needed structure",
		}
	}

	if !tokensBody.Success {
		return nil, &fiber.Error{
			Code:    tokensBody.StatusCode,
			Message: "Got an error from external API. See code",
		}
	}

	return tokensBody.Result.Prices, nil
}

func GetQuoteApi(query QuoteQuery) (*QuoteResultAPI, error) {
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(cconst.SwapAPIURL+
		"/quote?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&chainId=%d&gasPrice=%d",
		query.FromTokenAddress, query.ToTokenAddress, query.Amount,
		query.ChainId, query.GasPrice))

	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}

	var quoteBody QuoteResponseAPI

	err = json.Unmarshal(res.Body(), &quoteBody)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Bad JSON from external API or not successfully response",
		}
	}
	if !quoteBody.Success {
		return nil, &fiber.Error{
			Code:    quoteBody.StatusCode,
			Message: "Got error from external API. See status code",
		}
	}
	return &quoteBody.Result, nil
}

func GetSwapApi(query SwapQuery) (*SwapResultAPI, error) {
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(cconst.SwapAPIURL+
		"/swap?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&chainId=%d&slippage=%d&fromAddress=%s",
		query.FromTokenAddress, query.ToTokenAddress, query.Amount,
		query.ChainId, query.Slippage, query.FromAddress))

	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}

	var swapBody SwapResponseAPI

	err = json.Unmarshal(res.Body(), &swapBody)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Bad JSON from external API or not successfully response",
		}
	}
	if !swapBody.Success {
		return nil, &fiber.Error{
			Code:    swapBody.StatusCode,
			Message: "Got error from external API. See status code",
		}
	}
	return &swapBody.Result, nil
}

// Function
