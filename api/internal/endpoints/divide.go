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

func getTokenPrice(token models.TokenInPortfolio, tokensPrices map[string]string) (string, error) {
	tokenAmount, _, err := big.ParseFloat(token.Amount, 10, token.Decimals, big.AwayFromZero)
	if err != nil {
		return "", err
	}
	log.Println(tokenAmount.String())
	return "", nil
}

func sumTokensPriceInPortfolio(portfolio models.PortfolioResponse) (string, error) {
	tokensPrices, err := etc.GetTokensPrices(portfolio.ChainId)
	if err != nil {
		return "", err
	}
	nativeToken, err := etc.GetNativeTokenInfo(portfolio.ChainId)
	sum := decimal.NewFromInt(0)
	for _, token := range portfolio.Tokens {
		_price, success := big.NewInt(0).SetString(tokensPrices[token.Address], 10)
		if !success {
			return "", types.Error{Msg: "Not valid price type"}
		}
		realPrice := decimal.NewFromBigInt(_price, -int32(nativeToken.TokenDecimals))
		_amount, success := big.NewInt(0).SetString(token.Amount, 10)
		realAmount := decimal.NewFromBigInt(_amount, -int32(token.Decimals))
		sum = sum.Add(realPrice.Mul(realAmount))
	}

	return sum.String(), nil
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
	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(map[string]string{"sum": sum})
}
