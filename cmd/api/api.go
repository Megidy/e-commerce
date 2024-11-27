package api

import (
	"database/sql"
	"log"

	"github.com/Megidy/e-commerce/frontend/response"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/services/cart"
	"github.com/Megidy/e-commerce/services/information"
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
	informationStore := information.NewStore(s.db)

	TemplateHandler := response.NewTemplateHandler()
	router.RedirectFixedPath = true
	router.RedirectTrailingSlash = true

	authHandler := auth.NewJWT(userStore)
	managerHandler := auth.NewManager(userStore)

	userHandler := user.NewHandler(TemplateHandler, userStore)
	userHandler.RegisterRoutes(router, authHandler, managerHandler)

	productHandler := product.NewHandler(TemplateHandler, userStore, productStore)
	productHandler.RegisterRoutes(router, authHandler, managerHandler)

	cartHandler := cart.NewHandler(userStore, cartStore, productStore, TemplateHandler)
	cartHandler.RegisterRoutes(router, authHandler)

	orderHandler := order.NewHandler(TemplateHandler, userStore, orderStore, productStore, cartStore)
	orderHandler.RegisterRoutes(router, authHandler)

	informationHandler := information.NewHandler(TemplateHandler, userStore, informationStore)
	informationHandler.RegisterRoutes(router, authHandler, managerHandler)

	log.Println("Started Server on port :8080")
	return router.Run(s.addr)

}
