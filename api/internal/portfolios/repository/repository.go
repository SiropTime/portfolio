package token

import (
	"github.com/gofiber/fiber/v2"
	"go/types"
	"portfolioTask/api/internal/portfolios"
	"portfolioTask/api/internal/token"
	"portfolioTask/api/internal/token/repository"
	"portfolioTask/api/pkg/storage"
)

func CreatePortfolio(portfolio portfolios.InputPortfolio) (*portfolios.ProportionsResponsePortfolio, error) {
	conn, err := storage.CreateConnection()
	if err != nil {
		return nil, err
	}

	portfolioResult := conn.QueryRowx(`	
				INSERT INTO portfolios (chain_id, name) VALUES ($1, $2)
				RETURNING *;
	`, portfolio.ChainId, portfolio.Name)
	if err != nil {
		return nil, err
	}

	var portfolioDB portfolios.DBPortfolio

	if portfolioResult != nil {
		err = portfolioResult.StructScan(&portfolioDB)
	} else {
		return nil, types.Error{Msg: "Can't create connection to DB"}
	}
	err = repository.AddTokens(conn, portfolioDB.Id, portfolioDB.ChainId, portfolio.Tokens)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusNotAcceptable,
			Message: "Can't correctly insert tokens in portfolios",
		}
	}
	_p, err := ReadPortfolio(portfolioDB.Id)
	portfolioResponse, err := portfolios.CalculatePortfolioProportions(_p)
	if err != nil {
		return nil, &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "New portfolios wasn't created",
		}
	}
	return portfolioResponse, nil
}

func ReadPortfolio(id int) (portfolios.ResponsePortfolio, error) {
	conn, err := storage.CreateConnection()
	if err != nil {
		return portfolios.ResponsePortfolio{}, err
	}

	// Getting portfolios
	portfolioResult := conn.QueryRowx(`
		SELECT * FROM portfolios WHERE id = $1;
	`, id)
	var portfolioDB portfolios.DBPortfolio
	if portfolioResult != nil {
		err = portfolioResult.StructScan(&portfolioDB)
	} else {
		return portfolios.ResponsePortfolio{}, fiber.ErrNotFound
	}

	// Getting tokens inside the portfolios
	tokensResult, err := conn.Queryx(`
					SELECT amount, address, ticker, decimals FROM tokens
						WHERE portfolio_id = $1;
	`, portfolioDB.Id)

	var tokens []token.PortfoliosToken

	for tokensResult.Next() {
		var _token token.PortfoliosToken
		err = tokensResult.StructScan(&_token)
		tokens = append(tokens, _token)
	}

	// Unite this
	var portfolio portfolios.ResponsePortfolio
	portfolio.Id = portfolioDB.Id
	portfolio.ChainId = portfolioDB.ChainId
	portfolio.Name = portfolioDB.Name
	portfolio.Tokens = tokens

	return portfolio, nil
}

func ReadAllPortfolios() ([]portfolios.ResponsePortfolio, error) {
	conn, err := storage.CreateConnection()
	if err != nil {
		return nil, err
	}
	prePortfolios, err := conn.Queryx(`
		SELECT * FROM portfolios;
	`)

	var listPortfolios []portfolios.DBPortfolio

	for prePortfolios.Next() {
		_pdb := portfolios.DBPortfolio{}
		err = prePortfolios.StructScan(&_pdb)
		listPortfolios = append(listPortfolios, _pdb)
	}

	var resultList []portfolios.ResponsePortfolio
	for _, p := range listPortfolios {
		portfolio, err := ReadPortfolio(p.Id)
		if err != nil {
			continue
		}
		resultList = append(resultList, portfolio)
	}

	if resultList != nil {
		return resultList, nil
	}
	return []portfolios.ResponsePortfolio{}, types.Error{Msg: "Got empty portfolios, check if there is data in DB"}
}

func UpdatePortfolio(portfolioId int, portfolio portfolios.InputPortfolio) error {
	// PUT
	conn, err := storage.CreateConnection()
	if err != nil {
		return err
	}

	_, err = conn.Queryx(`
				      UPDATE portfolios
					  SET chain_id = $1, name = $2
					  WHERE id = $3
				      `, portfolio.ChainId, portfolio.Name, portfolioId)
	if err != nil {
		return err
	}
	err = repository.AddTokens(conn, portfolioId, portfolio.ChainId, portfolio.Tokens)
	if err != nil {
		return nil
	}
	return nil
}

func DeletePortfolio(id int) error {
	conn, err := storage.CreateConnection()
	if err != nil {
		return err
	}
	_, err = conn.Queryx(`
					  DELETE FROM portfolios
						WHERE id = $1;
					  `, id)
	if err != nil {
		return err
	}
	return nil
}

func AddNewTokens(portfolioId int, tokens token.InputTokens) error {
	// PATCH
	conn, err := storage.CreateConnection()
	if err != nil {
		return err
	}

	portfolio, err := ReadPortfolio(portfolioId)

	if err != nil {
		return err
	}
	tokensList := tokens.Tokens
	for _, tokenDetailed := range tokensList {
		tokenAPI, err := repository.GetTokenDetails(portfolio.ChainId, tokenDetailed.Address)
		if err != nil {
			return err
		}

		_, err = conn.Queryx(`
								UPDATE tokens SET
									  amount = $1
								  WHERE portfolio_id = $2 AND address = $3
									`, tokenDetailed.Amount, portfolioId, tokenDetailed.Address)

		if err != nil {
			_, e := conn.Queryx(`
					 INSERT INTO tokens VALUES
						(default, $1, $2, $3, $4, $5)
					 `, portfolioId, tokenDetailed.Amount, tokenDetailed.Address,
				tokenAPI.TokenTicker, tokenAPI.TokenDecimals)

			if e != nil {
				return e
			}
		}
	}

	return nil
}
