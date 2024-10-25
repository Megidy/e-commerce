package db

import (
	"database/sql"
	"log"

	"github.com/Megidy/e-commerce/config"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQlStorage() (*sql.DB, error) {
	dsn := GetDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = InitStorage(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDSN() string {
	config := *config.InitConfig()
	log.Println(config)
	return config.DBUser + ":" + config.DBPassword + "@" + config.DBProtocol + "(" + config.DB + ":" + config.DBPort + ")" + "/" + config.DBName
}

func InitStorage(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	log.Println("Success: DB is UP")
	return nil
}
