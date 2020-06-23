package moex

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func CreateTable() *sql.DB {
	var db *sql.DB
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		fmt.Println("open: ", err)
	}

	if err := DropTable(db); err != nil {
		fmt.Println("drop: ", err)
	}

	if _, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS securities(
		ticker varchar (100), 
		shortname varchar (100), 
		name varchar (100),
		latname varchar (100), 
		price real,
		capital double precision);
	CREATE TABLE IF NOT EXISTS prev_securities(
		ticker varchar (100), 
		price real);
	CREATE TABLE IF NOT EXISTS states(
		userid integer primary key, 
		state integer,
		listid integer,
		ticker varchar (100));
	CREATE TABLE IF NOT EXISTS lists(
		userid integer, 
		listid serial not null primary key,
		name varchar (100),
		state integer);
	CREATE TABLE IF NOT EXISTS list_items(
		listid integer, 
		ticker varchar (100),
		mode integer,
		lower real, 
		upper real);
		`); err != nil {
		fmt.Println("create:", err)
	}

	return db
}

func DropTable(db *sql.DB) error {
	if _, err := db.Exec(`
	DROP TABLE securities;
	DROP TABLE prev_securities;
	DROP TABLE states;
	DROP TABLE lists;
	DROP TABLE list_items;
	`); err != nil {
		return err
	}

	return nil
}

func SaveState(db *sql.DB, userid int, listid int, ticker string, state int) (err error) {
	res, err := db.Exec("UPDATE states SET state = $1, listid = $2, ticker = $3"+
		"WHERE userid = $4", state, listid, ticker, userid)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not update state in database: %v", err)
	}

	if aff, err := res.RowsAffected(); aff == 0 || err != nil {
		_, err = db.Exec("INSERT INTO states(userid, state, listid, ticker)"+
			"values($1, $2, $3, $4)", userid, state, listid, ticker)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("could not insert or replace state database entry: %v", err)
		}
	}

	return nil
}

func GetState(db *sql.DB, userid int) (state int, listid int, ticker string, err error) {
	row := db.QueryRow("SELECT state, listid, ticker FROM states WHERE userid = $1", userid)
	if err = row.Scan(&state, &listid, &ticker); err != nil {
		return state, listid, ticker, err
	}
	return state, listid, ticker, nil
}
