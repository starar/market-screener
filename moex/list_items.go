package moex

import (
	"database/sql"
	"fmt"
	"log"
)

type Item struct {
	Listid int
	Ticker string
	Mode   int
	Lower  float64
	Upper  float64
}

func InsertItem(db *sql.DB, listid int, ticker string, mode int) (err error) {
	_, err = db.Exec(`INSERT INTO list_items(listid, ticker, mode) VALUES($1, $2, $3);`,
		listid, ticker, mode)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not update state in database: %v", err)
	}
	return nil
}

func UpdateItem(db *sql.DB, listid int, ticker string, lower float64, upper float64) (err error) {
	if _, err := db.Exec(
		`UPDATE list_items SET lower = $1, upper = $2 WHERE listid = $3 AND ticker = $4;`,
		lower,
		upper,
		listid,
		ticker); err != nil {
		return err
	}
	return nil
}

func GetItem(db *sql.DB, listid int, ticker string) (mode int, lower float64, upper float64, err error) {
	row := db.QueryRow("SELECT mode, lower, upper FROM list_items WHERE listid = $1 AND ticker = $2",
		listid, ticker)
	if err := row.Scan(&mode, &lower, &upper); err != nil {
		fmt.Println("IsListNameFree:", err)
		return mode, lower, upper, err
	}
	return mode, lower, upper, nil
}

func GetItemsAll(db *sql.DB) (response []Item) {
	row, _ := db.Query("SELECT listid, ticker, mode, lower, upper FROM list_items WHERE upper IS NOT NULL;")
	defer row.Close()
	for row.Next() {
		var it Item
		if err := row.Scan(&it.Listid, &it.Ticker, &it.Mode, &it.Lower, &it.Upper); err != nil {
			log.Fatal(err)
		}
		response = append(response, it)
	}
	return
}

func GetItemsByList(db *sql.DB, listid int) (tickers []string, modes []int, lower []float64, upper []float64) {
	row, _ := db.Query("SELECT ticker, mode, lower, upper FROM list_items WHERE listid = $1;", listid)
	defer row.Close()
	for row.Next() {
		var ticker string
		var mode int
		var low, up float64
		if err := row.Scan(&ticker, &mode, &low, &up); err != nil {
			log.Fatal(err)
		}
		tickers = append(tickers, ticker)
		modes = append(modes, mode)
		lower = append(lower, low)
		upper = append(upper, up)
	}
	return
}

func IsItemTickerFree(db *sql.DB, listid int, ticker string) bool {
	row := db.QueryRow("SELECT ticker FROM list_items WHERE listid = $1 AND ticker = $2",
		listid, ticker)
	if err := row.Scan(&ticker); err != nil {
		fmt.Println("IsListNameFree:", err)
		return true
	}
	return false
}
