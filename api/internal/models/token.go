package models

type TokenPrice struct {
	Address string `json:"address"`
	Price   int64  `json:"price"`
}

type TokenInPortfolio struct {
	Ticker   string `json:"ticker"`
	Address  string `json:"address"`
	Decimals int16  `json:"decimals"`
	Amount   string `json:"amount"`
}

type TokenDB struct {
	portfolioId int    `db:"portfolio_id"`
	amount      string `db:"amount"`
	address     string `db:"address"`
	ticker      string `db:"ticker"`
	decimals    int16  `db:"decimals"`
}

type TokenRequestAPI struct {
	Success    bool       `json:"success"`
	StatusCode int        `json:"status_code"`
	Result     []TokenAPI `json:"result"`
}

type TokenAPI struct {
	ChainId              int      `json:"chain_id"`
	Tags                 []string `json:"tags"`
	TokenContractAddress string   `json:"token_contract_address"`
	TokenDecimals        int      `json:"token_decimals"`
	TokenImageUrl        string   `json:"token_image_url"`
	TokenName            string   `json:"token_name"`
	TokenTicker          string   `json:"token_ticker"`
}

type TokenPriceAPI struct {
	Success    bool              `json:"success"`
	StatusCode int               `json:"status_code"`
	Result     map[string]string `json:"result"`
}
