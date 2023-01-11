package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB(dsn string) (*sql.DB, error) {

	if !strings.Contains(dsn, "parseTime=true") {
		return nil, fmt.Errorf("dsn must include '?parseTime=true' ")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) ConnectToDB() (*sql.DB, error) {
	connection, err := OpenDB(app.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	app.Logger.Info("Connected to App DB")
	return connection, err
}
