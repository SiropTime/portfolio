package repositories

import (
	"awesomeProject/api/internal/models"
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
	err = firstInitialization(conn, models.Schema)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func firstInitialization(db *sqlx.DB, schema string) error {
	_, err := db.Exec(schema)

	return err
}
