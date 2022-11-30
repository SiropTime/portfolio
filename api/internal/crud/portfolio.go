package crud

import (
	"awesomeProject/api/internal/models"
	"awesomeProject/api/pkg/repositories"
	"github.com/jmoiron/sqlx"
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

func Read(id int) models.PortfolioResponse {
	conn, err := repositories.CreateConnection()
	if err != nil {
		log.Fatalln("Can't create connection with DB")
	}

	// Getting portfolio
	portfolioResult := conn.QueryRowx(`
		SELECT * FROM portfolios WHERE id = $1;
	`, id)
	var portfolioDB models.PortfolioDB
	if portfolioResult != nil {
		err = portfolioResult.StructScan(&portfolioDB)
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
	portfolio.Tokens = tokens

	return portfolio
}

func ReadAll() []models.PortfolioDB {
	return []models.PortfolioDB{}
}

func Update(db *sqlx.DB, portfolio models.PortfolioDB) error {
	return nil
}

func Delete(db *sqlx.DB, id int) error {
	return nil
}
