package models

// QuoteInput Client's input
type QuoteInput struct {
	PortfolioId int    `json:"portfolio_id"`
	GasPrice    int    `json:"gas_price"`
	TotalAmount string `json:"total_amount"`
}

// QuoteResult Common result from request
type QuoteResult struct {
	PortfolioId     int               `json:"portfolio_id"`
	ChainId         int               `json:"chain_id"`
	ConvertedTokens []CalculatedToken `json:"tokens"`
}

// CalculatedToken Converted token amount from quote request
type CalculatedToken struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}
