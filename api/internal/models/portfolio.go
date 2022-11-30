package models

type PortfolioDB struct {
	Id      int    `db:"id"`
	ChainId int    `db:"chain_id"`
	Name    string `db:"name"`
}

type PortfolioResponse struct {
	Id      int                `json:"id"`
	ChainId int                `json:"chain_id"`
	Name    string             `json:"name"`
	Tokens  []TokenInPortfolio `json:"tokens"`
	//Total   int                `json:"total"`
}

type PortfolioInput struct {
	ChainId int          `json:"chain_id"`
	Name    string       `json:"name"`
	Tokens  []TokenInput `json:"tokens"`
}

// Schema for future migrations
var Schema = `
	CREATE TABLE IF NOT EXISTS portfolios (
	    id SERIAL PRIMARY KEY,
	    chain_id INTEGER,
	    name TEXT
	);
	CREATE TABLE IF NOT EXISTS tokens_addreses (
	    id SERIAL PRIMARY KEY,
	    portfolio_id INT REFERENCES portfolios(id) ON DELETE CASCADE,
	    amount TEXT,
		address VARCHAR(48),
		ticker VARCHAR(16)                                
	);
`
