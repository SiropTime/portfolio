package models

type Token struct {
	Amount float32 `db:"amount"`
	Name   string  `db:"name"`
}

type PortfolioDB struct {
	UserId string  `db:"user_id"`
	Tokens []Token `db:"tokens"`
}

var Schema = `
	CREATE TABLE IF NOT EXISTS tokens (
		id SERIAL PRIMARY KEY,
		name VARCHAR(64),
		amount float
		);
	
	CREATE TABLE IF NOT EXISTS portfolios (
	    id SERIAL PRIMARY KEY,
	    user_id VARCHAR(64),
	    tokens integer[]

	)
`
