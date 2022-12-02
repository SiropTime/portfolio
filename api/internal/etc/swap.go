package etc

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"portfolioTask/api/internal/models"
	"portfolioTask/api/pkg/constants"
)

// File for interaction with API

func GetNativeTokenInfo(chainId int) (models.TokenAPI, error) {
	tokens, err := GetTokensAPI(chainId)
	if err != nil {
		return models.TokenAPI{}, err
	}
	for _, token := range tokens {
		if token.TokenContractAddress == constants.NativeAddress {
			return token, nil
		}
	}
	return models.TokenAPI{}, fiber.ErrNotFound
}

func GetTokensAPI(chainId int) ([]models.TokenAPI, error) {
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(constants.SwapAPIURL+"/tokens?chainId=%d", chainId))
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}
	var tokensBody models.TokenRequestAPI
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
	res, err := client.R().Get(fmt.Sprintf(constants.SwapAPIURL+"/prices?chainId=%d", chainId))
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}
	var tokensBody models.TokenPriceAPI

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

func GetQuoteApi(query models.QuoteQuery) (*models.QuoteResultAPI, error) {
	client := resty.New()
	res, err := client.R().Get(fmt.Sprintf(constants.SwapAPIURL+
		"/quote?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&chainId=%d&gasPrice=%d",
		query.FromTokenAddress, query.ToTokenAddress, query.Amount,
		query.ChainId, query.GasPrice))

	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to external API",
		}
	}

	var quoteBody models.QuoteResponseAPI

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

// Function
