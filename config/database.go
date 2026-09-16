package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDatabase() (*sql.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/db-sistem-analisis-ukom"

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