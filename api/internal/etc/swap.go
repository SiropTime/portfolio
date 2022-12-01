package etc

import (
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/constants"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// File for interaction with API

func GetChains() ([]models.Chain, error) {
	req, err := http.NewRequest(http.MethodGet,
		constants.SwapAPIURL+"/chains",
		nil,
	)

	if err != nil {
		log.Println("Can't form request")
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var chainsBody models.ChainsRequest
	err = json.Unmarshal(body, &chainsBody)
	if err != nil || !chainsBody.Success {
		return nil, err
	}
	return chainsBody.Result.Chains, nil
}

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
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(constants.SwapAPIURL+"/tokens?chainId=%d", chainId),
		nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var tokensBody models.TokenRequestAPI
	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &tokensBody)
	if err != nil || !tokensBody.Success {
		return nil, err
	}

	return tokensBody.Result, nil
}

func GetTokensPrices(chainId int) (map[string]string, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(constants.SwapAPIURL+"/prices?chainId=%d", chainId),
		nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var tokensBody models.TokenPriceAPI

	body, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &tokensBody)
	if err != nil || !tokensBody.Success {
		return nil, err
	}

	return tokensBody.Result.Prices, nil
}

func GetQuoteApi(query models.QuoteQuery) (*models.QuoteResultAPI, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(constants.SwapAPIURL+
			"/quote?fromTokenAddress=%s&toTokenAddress=%s&amount=%s&chainId=%d&gasPrice=%d",
			query.FromTokenAddress, query.ToTokenAddress, query.Amount,
			query.ChainId, query.GasPrice),
		nil)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't create request to API",
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't send request to API",
		}
	}
	defer res.Body.Close()
	var quoteBody models.QuoteResponseAPI
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't proceed body or not valid body of response",
		}
	}

	err = json.Unmarshal(body, &quoteBody)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "Bad JSON from external API or not successfully response",
		}
	}
	if !quoteBody.Success {
		return nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: string(body),
		}
	}
	return &quoteBody.Result, nil
}
