package main

import (
	"log"

	"github.com/Megidy/e-commerce/cmd/api"
	"github.com/Megidy/e-commerce/db"
)

func main() {
	db, err := db.NewMySQlStorage()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewApiServer(":8080", db)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}

}
