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

	conn.Exec(`	
				INSERT INTO portfolios (chain_id) VALUES (?)
	`, portfolio.ChainId)
	return nil
}

func Read(db *sqlx.DB, id int) models.PortfolioDB {
	return models.PortfolioDB{}
}

func ReadAll(db *sqlx.DB) []models.PortfolioDB {
	return []models.PortfolioDB{}
}

func Update(db *sqlx.DB, portfolio models.PortfolioDB) error {
	return nil
}

func Delete(db *sqlx.DB, id int) error {
	return nil
}
