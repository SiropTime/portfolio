package crud

import (
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/repositories"
	"github.com/jmoiron/sqlx"
	"go/types"
	"log"
)

func Create(portfolio models.PortfolioDB) error {
	conn, err := repositories.CreateConnection()
	if err != nil {
		log.Fatalln("Can't create connection with DB")

	}

	_, err = conn.Exec(`	
				INSERT INTO portfolios (chain_id) VALUES (?)
	`, portfolio.ChainId)
	if err != nil {
		return err
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

func Update(db *sqlx.DB, portfolio models.PortfolioDB) error {
	return nil
}

func Delete(db *sqlx.DB, id int) error {
	return nil
}
