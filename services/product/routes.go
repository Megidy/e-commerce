package product

import (
	"log"

	templates "github.com/Megidy/e-commerce/frontend/templates/product"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userStore    types.UserStore
	productStore types.ProductStore
}

func NewHandler(userStore types.UserStore, productStore types.ProductStore) *Handler {
	return &Handler{
		userStore:    userStore,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/products/accessories", h.GetAllAccessories)
	router.GET("/products/bicycles", h.GetAllBicycles)
}

func (h *Handler) GetAllAccessories(c *gin.Context) {
	accessories, err := h.productStore.GetAllAccessories()
	if err != nil {
		log.Println(err)
		return
	}
	templates.LoadAccessories(accessories).Render(c.Request.Context(), c.Writer)
}
func (h *Handler) GetAllBicycles(c *gin.Context) {
	bicycles, err := h.productStore.GetAllBicycles()
	if err != nil {
		log.Println(err)
		return
	}
	templates.LoadBicycles(bicycles).Render(c.Request.Context(), c.Writer)
}
