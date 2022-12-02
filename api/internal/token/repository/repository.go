package repository

import (
	"github.com/jmoiron/sqlx"
	"go/types"
	"portfolioTask/api/internal/token"
	"portfolioTask/api/pkg/clients/swapAPI"
)

func GetTokenDetails(chainId int, tokenAddress string) (swapAPI.TokenAPI, error) {
	tokensAPI, err := swapAPI.GetTokensAPI(chainId)
	if err != nil {
		return swapAPI.TokenAPI{}, err
	}
	for _, tokenDetailed := range tokensAPI {
		if tokenDetailed.TokenContractAddress == tokenAddress {
			return tokenDetailed, nil
		}
	}
	return swapAPI.TokenAPI{}, types.Error{Msg: "Couldn't find tokenDetailed with this address in current chain. Please, check one of the parameters"}

}

func AddTokens(conn *sqlx.DB, portfolioId int, chainId int, tokens []token.InputToken) error {
	for _, tokenNew := range tokens {
		var tokenDB token.DBToken
		tokenDB.PortfolioId = portfolioId
		tokenDB.Address, tokenDB.Amount = tokenNew.Address, tokenNew.Amount
		_t, err := GetTokenDetails(chainId, tokenDB.Address)
		if err != nil {
			return err
		}
		tokenDB.Decimals = _t.TokenDecimals
		tokenDB.Ticker = _t.TokenTicker

		_, err = conn.Queryx(`
						INSERT INTO tokens VALUES
						   (default, $1, $2, $3, $4, $5)
                         `,
			tokenDB.PortfolioId, tokenDB.Amount,
			tokenDB.Address, tokenDB.Ticker,
			tokenDB.Decimals)
		if err != nil {
			return err
		}
	}
	return nil
}
