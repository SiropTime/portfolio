package endpoints

import (
	"awesomeProject/api/internal/crud"
	"awesomeProject/api/internal/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetPortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	portRes, err := crud.Read(pId)
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portRes)
}

func GetAllPortfolios(c *fiber.Ctx) error {
	portfolios, err := crud.ReadAll()
	if err != nil {
		return err
	}

	return json.NewEncoder(c.Type("json", "utf-8").Response().BodyWriter()).Encode(portfolios)
}

func PostPortfolio(c *fiber.Ctx) error {
	portfolio := new(models.PortfolioInput)
	if err := c.BodyParser(portfolio); err != nil {
		c.Status(503)
	}
	err := crud.Create(*portfolio)
	if err != nil {
		return err
	}
	return nil
}

func DeletePortfolio(c *fiber.Ctx) error {
	pId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = crud.Delete(pId)
	if err != nil {
		return err
	}
	return nil
}

func PrepareTransaction(c *fiber.Ctx) {

}
