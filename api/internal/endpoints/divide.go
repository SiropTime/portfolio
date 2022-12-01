package endpoints

import (
	"awesomeProject/api/internal/crud"
	"awesomeProject/api/internal/etc"
	"awesomeProject/api/internal/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go/types"
	"log"
	"math/big"
	"strconv"
)

func calculateRealAmount(token models.TokenInPortfolio) {

}

func getTokenPrice(token models.TokenInPortfolio, tokensPrices map[string]string) (string, error) {
	tokenAmount, _, err := big.ParseFloat(token.Amount, 10, token.Decimals, big.AwayFromZero)
	if err != nil {
		return "", err
	}
	log.Println(tokenAmount.String())
	return "", nil
}

func sumTokensPriceInPortfolio(portfolio models.PortfolioResponse) (*big.Float, error) {
	tokensPrices, err := etc.GetTokensPrices(portfolio.ChainId)
	if err != nil {
		return nil, err
	}
	nativeToken, err := etc.GetNativeTokenInfo(portfolio.ChainId)
	sum := new(big.Float).SetMantExp(big.NewFloat(0), -int(nativeToken.TokenDecimals))
	for _, token := range portfolio.Tokens {
		_price, success := new(big.Float).SetString(tokensPrices[token.Address])
		realPrice := new(big.Float).SetMantExp(_price, -int(nativeToken.TokenDecimals))
		if !success {
			return nil, types.Error{Msg: "There's no valid price string in data"}
		}
		_amount, success := new(big.Float).SetString(token.Amount)
		realAmount := new(big.Float).SetMantExp(_amount, -int(token.Decimals))
		_temp := new(big.Float).SetInt64(0).SetPrec(nativeToken.TokenDecimals)
		_temp.Mul(realPrice, realAmount)
		sum = big.NewFloat(0.0).SetPrec(nativeToken.TokenDecimals).Add(sum, _temp)

	}

	return sum, nil
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

	sum, err := sumTokensPriceInPortfolio(portfolio)
	if err != nil {
		return err
	}
	//c.SendString(sum.String())
	//return nil
	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(map[string]*big.Float{"sum": sum})
}
