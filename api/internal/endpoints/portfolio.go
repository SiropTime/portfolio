package endpoints

import (
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/constants"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func getChains() models.ChainsList {
	req, err := http.NewRequest(http.MethodGet,
		constants.SwapAPIURL+"/chains",
		nil,
	)

	if err != nil {
		log.Println("Can't get to SwapAPI")
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Can't get to SwapAPI")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var chainsBody models.ChainsRequest
	err = json.Unmarshal(body, &chainsBody)
	if err != nil || !chainsBody.Success {
		log.Println(err)
	}
	return chainsBody.Result
}

func GetPortfolio(c *fiber.Ctx) error {
	body := getChains()
	err := c.SendString(strconv.Itoa(body.Chains[0].ChainId))
	if err != nil {
		return err
	}
	return nil
}

func PrepareTransaction(c *fiber.Ctx) {

}
