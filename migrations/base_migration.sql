create table portfolios
(
    id serial primary key,
    chain_id integer,
    name     text
);
create table if not exists tokens
(
    id serial primary key,
    portfolio_id integer
        constraint tokens_addreses_portfolio_id_fkey
            references portfolios
            on delete cascade,
    amount       text,
    address      varchar(48),
    ticker       varchar(16),
    decimals     integer
);