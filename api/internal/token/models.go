package token

import (
	"github.com/shopspring/decimal"
)

type PriceToken struct {
	Address string `json:"address"`
	Price   int64  `json:"price"`
}

type InputToken struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type InputTokens struct {
	Tokens []InputToken `json:"tokens"`
}

type PortfoliosToken struct {
	Ticker   string `json:"ticker"`
	Address  string `json:"address"`
	Decimals uint   `json:"decimals"`
	Amount   string `json:"amount"`
}

type DBToken struct {
	PortfolioId int    `db:"portfolio_id"`
	Amount      string `db:"amount"`
	Address     string `db:"address"`
	Ticker      string `db:"ticker"`
	Decimals    uint   `db:"decimals"`
}

type RealDataToken struct {
	Ticker     string
	Address    string
	TotalPrice decimal.Decimal
}

type ProportionsToken struct {
	Ticker     string          `json:"ticker"`
	Address    string          `json:"address"`
	Proportion decimal.Decimal `json:"proportion"`
}

type PriceResult struct {
	LastUpdate string            `json:"last_update"`
	Prices     map[string]string `json:"prices"`
}

type QuoteToken struct {
	FinalAmount  string `json:"final_price"`
	NativePrice  string `json:"native_price"`
	EstimatedGas int    `json:"estimated_gas"`
	Address      string `json:"address"`
	Ticker       string `json:"ticker"`
}

// CalculatedToken Converted token amount from quote request
type CalculatedToken struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}
