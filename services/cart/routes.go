package cart

import (
	"log"
	"net/http"

	templates "github.com/Megidy/e-commerce/frontend/templates/cart"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userStore    types.UserStore
	cartStore    types.CartStore
	productStore types.ProductStore
}

func NewHandler(userStore types.UserStore, cartStore types.CartStore, productStore types.ProductStore) *Handler {
	return &Handler{userStore: userStore, cartStore: cartStore, productStore: productStore}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/cart", auth.NewHandler(h.userStore).WithJWT, h.GetCart)
}

func (h *Handler) GetCart(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(types.User)
	log.Println(user)
	cart, err := h.cartStore.GetCart(user.ID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	accessories, bicycles, err := h.productStore.GetAllProducts(cart)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
	}
	templates.LoadCart(bicycles, accessories).Render(c.Request.Context(), c.Writer)

}
