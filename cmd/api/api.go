package api

import "database/sql"

type APIServer struct {
	addr string
	db   *sql.DB
}
