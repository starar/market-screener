package moex

import "os"

var (
	//dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	dbInfo  = os.Getenv("DATABASE_URL")
	moexURL = "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities.xml"
	moex    = "https://iss.moex.com/iss/history/engines/stock/markets/shares/boards/TQBR/securities?date="
)

var (
	host     = "127.0.0.1"
	port     = "5432"
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
	sslmode  = "disable"
)
