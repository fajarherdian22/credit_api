package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConDB(config string) *sql.DB {
	db, err := sql.Open("mysql", config)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
