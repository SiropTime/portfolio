package models

type PortfolioDB struct {
	Id      int `db:"id"`
	ChainId int `db:"chain_id"`
}

type PortfolioResponse struct {
	Id      int                `json:"id"`
	ChainId int                `json:"chain_id"`
	Tokens  []TokenInPortfolio `json:"tokens"`
	//Total   int                `json:"total"`
}

type TokenInChain struct {
	Address string `json:"address"`
	Price   int64  `json:"price"`
}

type TokenInPortfolio struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
	Price   int64  `json:"price"`
	Amount  int64  `json:"Amount"`
}

type TokenDB struct {
	portfolioId int    `db:"portfolio_id"`
	amount      int64  `db:"amount"`
	address     string `db:"address"`
	symbol      string `db:"symbol"`
	price       int64  `db:"price"`
}

// Schema for future migrations
var Schema = `
	CREATE TABLE IF NOT EXISTS portfolios (
	    id SERIAL PRIMARY KEY,
	    chain_id INTEGER
	);
	CREATE TABLE IF NOT EXISTS tokens_addreses (
	    id SERIAL PRIMARY KEY,
	    portfolio_id INT REFERENCES portfolios(id) ON DELETE CASCADE,
	    amount BIGINT,
		address VARCHAR(48),
		symbol VARCHAR(16)                                
	);
`
