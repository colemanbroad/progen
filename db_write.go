package main

import (
	"database/sql"
	// "fmt"
	// "log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/marcboeker/go-duckdb"
	"golang.org/x/exp/rand"
)

// Logging History Datastructures
type Reward_power_of_two struct {
	Value  float32 `db:"value"`
	Reward float32 `db:"reward"`
	Time   int     `db:"time"`
}

var (
	history_power_of_two []Reward_power_of_two
)

func init_history() {
	history_power_of_two = make([]Reward_power_of_two, 0, 1_000)
}

// func Init_campaign() {
// 	init_library()
// 	init_history()
// 	init_maphistory()
// 	init_reward()
// }

func ConnectSqlite(dbname string) *sql.DB {
	driver := "sqlite3"
	db, err := sql.Open(driver, dbname)
	check(err)
	// if err != nil {
	// 	return db, fmt.Errorf("Error opening database: %v", err)
	// }
	err = db.Ping()
	check(err)
	// if err != nil {
	// 	ErrorLog.Fatal(err)
	// }
	return db
}

func Create_tables(dbname string) {
	db := ConnectSqlite(dbname)
	var s string
	var err error

	s = `create table if not exists history_power_of_two (value real, reward real, time int, campaign_id string)`
	_, err = db.Exec(s)
	check(err)

	s = `create table if not exists program_history (prog string, reward real, time int, campaign_id string)`
	_, err = db.Exec(s)
	check(err)

	s = `create table if not exists campaigns (campaign_id string, dtime datetime, n_iter int, ltype int)`
	_, err = db.Exec(s)
	check(err)

	db.Close()
}

// func Save_tables_or_fail(db *sql.DB) {
// 	err := save_tables_old(db)
// 	if err != nil {
// 		ErrorLog.Fatalf("Couldn't save tables: %v\n", err)
// 	}
// }

// func Save_tables_dbname(name string) {
// 	db := ConnectSqlite(name)
// 	err := save_tables_old(db)
// 	check(err)
// }

// func save_tables_old(db *sql.DB) error {
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	// Prepare the SQL insert statement
// 	stmt, err := tx.Prepare("insert into history_power_of_two (value, reward, time) values (?,?,?)")
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	for _, vr := range history_power_of_two {
// 		_, err = stmt.Exec(vr.Value, vr.Reward, vr.Time)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Commit the transaction
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
