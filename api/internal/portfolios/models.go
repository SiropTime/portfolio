package portfolios

import (
	"portfolioTask/api/internal/token"
)

// QuoteInput Client's input
type QuoteInput struct {
	PortfolioId int    `json:"portfolio_id"`
	GasPrice    int    `json:"gas_price"`
	TotalAmount string `json:"total_amount"`
}

// QuoteResult Common result from request
type QuoteResult struct {
	PortfolioId     int                     `json:"portfolio_id"`
	ChainId         int                     `json:"chain_id"`
	ConvertedTokens []token.CalculatedToken `json:"tokens"`
}

type DBPortfolio struct {
	Id      int    `db:"id"`
	ChainId int    `db:"chain_id"`
	Name    string `db:"name"`
}

type ResponsePortfolio struct {
	Id      int                     `json:"id"`
	ChainId int                     `json:"chain_id"`
	Name    string                  `json:"name"`
	Tokens  []token.PortfoliosToken `json:"tokens"`
	//Total   int                `json:"total"`
}

type InputPortfolio struct {
	ChainId int                `json:"chain_id"`
	Name    string             `json:"name"`
	Tokens  []token.InputToken `json:"tokens"`
}

type ProportionsResponsePortfolio struct {
	Id                int    `json:"id"`
	ChainId           int    `json:"chain_id"`
	Name              string `json:"name"`
	TokensProportions []token.ProportionsToken
}

type AfterQuotePortfolio struct {
	Id      int                `json:"id"`
	ChainId int                `json:"chain_id"`
	Name    string             `json:"name"`
	Tokens  []token.QuoteToken `json:"tokens"`
}
