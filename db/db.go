package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewPostgreSQLStorage(connectionSrt string) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@/productdb")
	if err != nil {
		panic(err)
	}
	return db, nil
}
