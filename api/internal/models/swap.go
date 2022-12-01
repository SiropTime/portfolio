package models

type Chain struct {
	ChainId       int    `json:"chain_id"`
	ChainImageUrl string `json:"chain_image_url"`
	ChainName     string `json:"chain_name"`
	CoingeckoId   string `json:"coingecko_id"`
	CurrencySym   string `json:"currency_sym"`
	ExplorerUrl   string `json:"explorer_url"`
	RpcUrl        string `json:"rpc_url"`
}

type ChainsList struct {
	Chains []Chain `json:"chains"`
}

type ChainsRequest struct {
	Success    bool       `json:"success"`
	StatusCode int        `json:"status_code"`
	Result     ChainsList `json:"result"`
}

type AnyRequest struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status_code"`
	Result     any  `json:"result"`
}
