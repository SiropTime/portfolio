package repositories

import (
	"awesomeProject/api/pkg/constants"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func CreateConnection() (*sqlx.DB, error) {
	conn, err := sqlx.Connect(
		constants.PostgresDriver,
		constants.PostgresUrl)
	if err != nil {
		log.Fatalln(err)
	}
	return conn, err
}

func FirstInitialization(db *sqlx.DB, schema string) error {
	_, err := db.Exec(schema)

	return err
}
