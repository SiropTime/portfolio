package constants

import (
	"fmt"
	"os"
)

const PostgresDriver = "postgres"

var PostgresUrl = fmt.Sprintf("host=postgresql port=5432 user=%s password=%s dbname=%s sslmode=disable",
	os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

const SwapAPIURL = "https://amaze.finance/api/swap/"
const NativeAddress = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
