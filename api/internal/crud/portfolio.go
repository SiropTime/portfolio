package crud

import (
	"awesomeProject/api/internal/models"
	"github.com/jmoiron/sqlx"
)

func Create(db *sqlx.DB, portfolio models.PortfolioDB) error {
	db.Exec(`	
				INSERT INTO portfolios () VALUES (??)
	`)
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
