package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/marcboeker/go-duckdb"
)

func Save_tables_fast(dbname string, params GPParams) error {

	connector, err := duckdb.NewConnector(dbname, nil)
	if err != nil {
		return err
	}
	defer connector.Close()

	conn, err := connector.Connect(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	// Retrieve appender from connection (note that you have to create the table 'test' beforehand).

	campaign_id := generateRandomString(16)
	// var appender *duckdb.Appender

	err = save_hist(conn, campaign_id)
	if err != nil {
		return err
	}
	err = save_proghist(conn, campaign_id)
	if err != nil {
		return err
	}

	stmt, err := conn.Prepare("insert into campaigns values (?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec([]driver.Value{campaign_id, time.Now(), params.N_rounds, params.Ltype})
	if err != nil {
		return err
	}
	stmt.Close()

	fmt.Println("Saved tables with campaign_id " + campaign_id)

	return nil
}

func save_hist(conn driver.Conn, campaign_id string) error {
	appender, err := duckdb.NewAppenderFromConn(conn, "", "history_power_of_two")
	if err != nil {
		return err
	}
	for _, row := range history_power_of_two {
		err = appender.AppendRow(row.Value, row.Reward, int32(row.Time), campaign_id)
		if err != nil {
			return err
		}
	}
	err = appender.Flush()
	if err != nil {
		return err
	}
	appender.Close()
	return nil
}

func save_proghist(conn driver.Conn, campaign_id string) error {
	appender, err := duckdb.NewAppenderFromConn(conn, "", "program_history")
	if err != nil {
		return err
	}
	for _, row := range program_history {
		err = appender.AppendRow(row.Prog, row.reward, int32(row.time), campaign_id)
		if err != nil {
			return err
		}
	}
	err = appender.Flush()
	if err != nil {
		return err
	}
	appender.Close()
	return nil
}
