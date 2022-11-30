package etc

import (
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/constants"
	"encoding/json"
	"fmt"
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

	return tokensBody.Result, nil
}
