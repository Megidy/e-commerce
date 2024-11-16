package api

import (
	"database/sql"
	"log"

	"github.com/Megidy/e-commerce/frontend/response"
	"github.com/Megidy/e-commerce/services/cart"
	"github.com/Megidy/e-commerce/services/order"
	"github.com/Megidy/e-commerce/services/product"
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
	router.Static("/static", "frontend/static")

	userStore := user.NewStore(s.db)
	productStore := product.NewStore(s.db)
	cartStore := cart.NewStore(s.db)
	orderStore := order.NewStore(s.db)
	NewResponseHandler := response.NewTemplateHandler()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	userHandler := user.NewHandler(NewResponseHandler, userStore)
	userHandler.RegisterRoutes(router)

	productHandler := product.NewHandler(userStore, productStore)
	productHandler.RegisterRoutes(router)

	cartHandler := cart.NewHandler(userStore, cartStore, productStore, NewResponseHandler)
	cartHandler.RegisterRoutes(router)

	orderHandler := order.NewHandler(userStore, orderStore, productStore, cartStore)
	orderHandler.RegisterRoutes(router)

	log.Println("Started Server on port :8080")
	return router.Run(s.addr)

}
