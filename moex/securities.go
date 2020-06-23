package moex

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Table struct {
	Data []struct {
		Id      string     `xml:"id,attr"`
		Results []Security `xml:"rows>row"`
	} `xml:"data"`
}

type Security struct {
	Ticker    string  `xml:"SECID,attr"`
	ShortName string  `xml:"SHORTNAME,attr"`
	Name      string  `xml:"SECNAME,attr"`
	LatName   string  `xml:"LATNAME,attr"`
	Price     float64 `xml:"PREVWAPRICE,attr"`
	Capital   float64 `xml:"ISSUESIZE,attr"`
}

func InsertSecurity(db *sql.DB, sec Security) error {
	sec.Capital *= sec.Price
	result, err := db.Exec(
		`UPDATE securities SET price = $1, capital = $2 WHERE ticker = $3;`,
		sec.Price,
		sec.Capital,
		sec.Ticker)

	if err != nil {
		fmt.Println("InsertSecurity:", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("InsertSecurity:", err)
		return err
	}
	if affected == 0 {
		if _, err := db.Exec(
			`INSERT INTO securities(ticker, shortname, name, latname, price, capital) VALUES($1, $2, $3, $4, $5, $6);`,
			sec.Ticker,
			sec.ShortName,
			sec.Name,
			sec.LatName,
			sec.Price,
			sec.Capital); err != nil {
			fmt.Println("InsertSecurity:", err)
			return err
		}
	}
	return nil
}

func UpdateSecurities(db *sql.DB) {
	fmt.Println("Updating...")

	response, err := http.Get(moexURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	table, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var result Table
	xml.Unmarshal([]byte(table), &result)

	for _, data := range result.Data {
		if data.Id != "securities" {
			continue
		}
		for _, sec := range data.Results {
			InsertSecurity(db, sec)
		}
	}
}

func GetTicker(db *sql.DB, ticker string) (price float64, capital float64) {
	row := db.QueryRow("SELECT price, capital FROM securities WHERE ticker = $1;", ticker)
	if err := row.Scan(&price, &capital); err != nil {
		log.Fatal(err)
	}
	return price, capital
}

func GetTickersAll(db *sql.DB) (response []Security) {
	row, _ := db.Query("SELECT shortname, ticker, price, capital FROM securities;")
	defer row.Close()
	for row.Next() {
		var sec Security
		if err := row.Scan(&sec.ShortName, &sec.Ticker, &sec.Price, &sec.Capital); err != nil {
			log.Fatal(err)
		}
		response = append(response, sec)
	}
	return
}

func FindCompany(db *sql.DB, keyWord string) []Security {
	row, _ := db.Query(
		"SELECT ticker, shortname, name, latname FROM securities WHERE ticker ~ '(?i)" + keyWord + "'" +
			"OR shortname ~ '(?i)" + keyWord + "'" +
			"OR name ~ '(?i)" + keyWord + "'" +
			"OR latname ~ '(?i)" + keyWord + "';")
	defer row.Close()

	var response []Security
	for row.Next() {
		var sec Security
		if err := row.Scan(&sec.Ticker, &sec.ShortName, &sec.Name, &sec.LatName); err != nil {
			log.Fatal(err)
		}
		response = append(response, sec)
	}

	return response
}

func IsTickerExist(db *sql.DB, ticker string) bool {
	row := db.QueryRow("SELECT ticker FROM securities WHERE ticker = $1", ticker)
	if err := row.Scan(&ticker); err != nil {
		return false
	}
	return true
}
