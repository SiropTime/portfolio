package models

type PortfolioDB struct {
	ChainId int `db:"chain_id"`
}

type PortfolioResponse struct {
	ChainId int `json:"chain_id"`
	Tokens  []TokenInChain
}

type TokenInChain struct {
	Address string `json:"address"`
	Price   int64  `json:"price"`
}

var Schema = `
	CREATE TABLE IF NOT EXISTS portfolios (
	    id SERIAL PRIMARY KEY,
	    user_id VARCHAR(64),
	    chain_id INTEGER
	);
`
