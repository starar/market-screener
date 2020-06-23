package moex

import (
	"database/sql"
	"fmt"
	"log"
)

func InsertList(db *sql.DB, userid int, name string, state int) (err error) {
	_, err = db.Exec(`INSERT INTO lists(userid, name, state) VALUES($1, $2, $3);`, userid, name, state)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not update state in database: %v", err)
	}
	return nil
}

func UpdateList(db *sql.DB, userid int, listid int, state int) error {
	if _, err := db.Exec(
		`UPDATE lists SET state = $1 WHERE listid = $2 AND userid = $3;`,
		state,
		listid,
		userid); err != nil {
		fmt.Println("UpdateList:", err)
		return err
	}
	return nil
}

func UpdateListsAll(db *sql.DB, userid int, state int) error {
	if _, err := db.Exec(
		`UPDATE lists SET state = $1 WHERE userid = $2;`,
		state, userid); err != nil {
		fmt.Println("StopTable:", err)
		return err
	}
	return nil
}

func GetListsAll(db *sql.DB, userid int, state int) (listids []int, names []string, states []int) {
	row, _ := db.Query("SELECT listid, name, state FROM lists WHERE userid = $1 AND state != $2;",
		userid, state)
	defer row.Close()
	for row.Next() {
		var name string
		var listid, state int
		if err := row.Scan(&listid, &name, &state); err != nil {
			log.Fatal(err)
		}
		listids = append(listids, listid)
		names = append(names, name)
		states = append(states, state)
	}
	return
}

func GetActiveLists(db *sql.DB, state int) (listids []int, userids []int, names []string) {
	row, err := db.Query("SELECT listid, userid, name FROM lists WHERE state = $1;", state)
	if err != nil {
		fmt.Println("active:", err)
		return
	}
	defer row.Close()
	for row.Next() {
		var listid, userid int
		var name string
		if err := row.Scan(&listid, &userid, &name); err != nil {
			log.Fatal(err)
		}
		listids = append(listids, listid)
		userids = append(userids, userid)
		names = append(names, name)
	}
	return
}

func GetListByName(db *sql.DB, userid int, name string) (listid int, err error) {
	row := db.QueryRow("SELECT listid FROM lists WHERE userid = $1 AND name = $2", userid, name)
	if err := row.Scan(&listid); err != nil {
		fmt.Println("GetListByName:", err)
		return 0, err
	}
	return listid, nil
}

func GetListState(db *sql.DB, userid int, listid int) (state int, err error) {
	row := db.QueryRow("SELECT state FROM lists WHERE userid = $1 AND listid = $2",
		userid, listid)
	if err = row.Scan(&state); err != nil {
		fmt.Println("GetListState:", err)
		return 0, err
	}
	return state, nil
}

func IsListNameFree(db *sql.DB, userid int, name string) bool {
	row := db.QueryRow("SELECT name FROM lists WHERE userid = $1 AND name = $2", userid, name)
	if err := row.Scan(&name); err != nil {
		fmt.Println("IsListNameFree:", err)
		return true
	}
	return false
}
