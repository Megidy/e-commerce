package api

import (
	"database/sql"
	"log"

	"github.com/Megidy/e-commerce/frontend/response"
	"github.com/Megidy/e-commerce/services/user"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {

	router := gin.Default()
	// router.LoadHTMLGlob("./frontend/templates/*.html")
	NewResponseHandler := response.NewTemplateHandler()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(NewResponseHandler, userStore)
	userHandler.RegisterRoutes(router)
	log.Println("Started Server on port :8080")
	return router.Run(s.addr)

}
