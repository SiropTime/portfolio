package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"portfolioTask/api/internal/cconst"
)

func CreateConnection() (*sqlx.DB, error) {
	conn, err := sqlx.Connect(
		cconst.PostgresDriver,
		cconst.PostgresUrl)
	if err != nil {
		log.Fatalln(err)
	}
	return conn, err
}
