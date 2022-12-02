package swapAPI

import (
	"github.com/shopspring/decimal"
	"portfolioTask/api/internal/token"
)

type TokenRequestAPI struct {
	Success    bool       `json:"success"`
	StatusCode int        `json:"status_code"`
	Result     []TokenAPI `json:"result"`
}

type TokenAPI struct {
	ChainId              int      `json:"chain_id"`
	Tags                 []string `json:"tags"`
	TokenContractAddress string   `json:"token_contract_address"`
	TokenDecimals        uint     `json:"token_decimals"`
	TokenImageUrl        string   `json:"token_image_url"`
	TokenName            string   `json:"token_name"`
	TokenTicker          string   `json:"token_ticker"`
}

type TokenPriceAPI struct {
	Success    bool              `json:"success"`
	StatusCode int               `json:"status_code"`
	Result     token.PriceResult `json:"result"`
}

type QuoteQuery struct {
	FromTokenAddress string
	ToTokenAddress   string
	Amount           string
	ChainId          int
	GasPrice         int
}

type QuoteResponseAPI struct {
	Success    bool           `json:"success"`
	StatusCode int            `json:"status_code"`
	Result     QuoteResultAPI `json:"result"`
}

type QuoteResultAPI struct {
	EstimatedGas    int    `json:"estimated_gas"`
	FromTokenAmount string `json:"from_token_amount"`
	ToTokenAmount   string `json:"to_token_amount"`
}

type TokensInfo struct {
	Total       decimal.Decimal       `json:"total"`
	NativeToken TokenAPI              `json:"native_token"`
	Tokens      []token.RealDataToken `json:"tokens"`
}
