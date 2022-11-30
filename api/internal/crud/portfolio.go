package crud

import (
	"awesomeProject/api/internal/etc"
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/repositories"
	"go/types"
)

func getTokenDetails(chainId int, tokenAddress string) (models.TokenAPI, error) {
	tokensAPI, err := etc.GetTokensAPI(chainId)
	if err != nil {
		return models.TokenAPI{}, err
	}
	for _, token := range tokensAPI {
		if token.TokenContractAddress == tokenAddress {
			return token, nil
		}
	}
	return models.TokenAPI{}, types.Error{Msg: "Couldn't find token with this address in current chain. Please, check one of the parameters"}

}

func Create(portfolio models.PortfolioInput) error {
	conn, err := repositories.CreateConnection()
	if err != nil {
		return err
	}

	portfolioResult := conn.QueryRowx(`	
				INSERT INTO portfolios (chain_id, name) VALUES ($1, $2)
				RETURNING *;
	`, portfolio.ChainId, portfolio.Name)
	if err != nil {
		return err
	}

	var portfolioDB models.PortfolioDB

	if portfolioResult != nil {
		err = portfolioResult.StructScan(&portfolioDB)
	} else {
		return types.Error{Msg: "Can't create connection to DB"}
	}

	for _, token := range portfolio.Tokens {
		var tokenDB models.TokenDB
		tokenDB.PortfolioId = portfolioDB.Id
		tokenDB.Address, tokenDB.Amount = token.Address, token.Amount
		_t, err := getTokenDetails(portfolioDB.ChainId, tokenDB.Address)
		if err != nil {
			return err
		}
		tokenDB.Decimals = _t.TokenDecimals
		tokenDB.Ticker = _t.TokenTicker

		_, err = conn.Queryx(`
						INSERT INTO tokens VALUES
						   (default, $1, $2, $3, $4, $5)
                         `,
			tokenDB.PortfolioId, tokenDB.Amount,
			tokenDB.Address, tokenDB.Ticker,
			tokenDB.Decimals)
		if err != nil {
			return err
		}
	}
	return nil
}

func Read(id int) (models.PortfolioResponse, error) {
	conn, err := repositories.CreateConnection()
	if err != nil {
		return models.PortfolioResponse{}, err
	}

	// Getting portfolio
	portfolioResult := conn.QueryRowx(`
		SELECT * FROM portfolios WHERE id = $1;
	`, id)
	var portfolioDB models.PortfolioDB
	if portfolioResult != nil {
		err = portfolioResult.StructScan(&portfolioDB)
	} else {
		return models.PortfolioResponse{}, types.Error{Msg: "Can't create connection to DB"}
	}

	// Getting tokens inside the portfolio
	tokensResult, err := conn.Queryx(`
					SELECT amount, address, ticker, decimals FROM tokens
						WHERE portfolio_id = $1;
	`, portfolioDB.Id)

	var tokens []models.TokenInPortfolio

	for tokensResult.Next() {
		var token models.TokenInPortfolio
		err = tokensResult.StructScan(&token)
		tokens = append(tokens, token)
	}

	// Unite this
	var portfolio models.PortfolioResponse
	portfolio.Id = portfolioDB.Id
	portfolio.ChainId = portfolioDB.ChainId
	portfolio.Name = portfolioDB.Name
	portfolio.Tokens = tokens

	return portfolio, nil
}

func ReadAll() ([]models.PortfolioResponse, error) {
	conn, err := repositories.CreateConnection()
	if err != nil {
		return nil, err
	}
	prePortfolios, err := conn.Queryx(`
		SELECT * FROM portfolios;
	`)

	var listPortfolios []models.PortfolioDB

	for prePortfolios.Next() {
		_pdb := models.PortfolioDB{}
		err = prePortfolios.StructScan(&_pdb)
		listPortfolios = append(listPortfolios, _pdb)
	}

	var resultList []models.PortfolioResponse
	for _, p := range listPortfolios {
		portfolio, err := Read(p.Id)
		if err != nil {
			continue
		}
		resultList = append(resultList, portfolio)
	}

	if resultList != nil {
		return resultList, nil
	}
	return []models.PortfolioResponse{}, types.Error{Msg: "Got empty portfolio, check if there is data in DB"}
}

func Update(portfolio models.PortfolioInput) error {
	// PUT
	_, err := repositories.CreateConnection()
	if err != nil {
		return err
	}

	return nil
}

func AddNewToken(token models.TokenInput) error {
	// PATCH
	_, err := repositories.CreateConnection()
	if err != nil {
		return err
	}
	return nil
}

func Delete(id int) error {
	conn, err := repositories.CreateConnection()
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
