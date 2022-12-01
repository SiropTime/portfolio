package repositories

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"portfolioTask/api/pkg/constants"
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

func FirstInitialization(schema string) error {
	conn, err := CreateConnection()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't connect to database"}
	}
	_, err = conn.Exec(schema)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: "Can't accept migrations"}
	}

	return err
}
