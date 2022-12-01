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

type PortfolioProportionsResponse struct {
	Id                int    `json:"id"`
	ChainId           int    `json:"chain_id"`
	Name              string `json:"name"`
	TokensProportions []TokenProportions
}

type PortfolioAfterQuote struct {
	Id      int          `json:"id"`
	ChainId int          `json:"chain_id"`
	Name    string       `json:"name"`
	Tokens  []TokenQuote `json:"tokens"`
}

// Schema for future migrations
var Schema = `
	create table portfolios
	(
    id       serial
        primary key,
    chain_id integer,
    name     text
	);
	create table if not exists tokens
	(
    id           integer default nextval('tokens_addreses_id_seq'::regclass) not null
        constraint tokens_addreses_pkey
            primary key,
    portfolio_id integer
        constraint tokens_addreses_portfolio_id_fkey
            references portfolios
            on delete cascade,
    amount       text,
    address      varchar(48),
    ticker       varchar(16),
    decimals     integer
	);

`
