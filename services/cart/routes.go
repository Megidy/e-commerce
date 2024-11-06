package cart

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	templates "github.com/Megidy/e-commerce/frontend/templates/cart"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userStore    types.UserStore
	cartStore    types.CartStore
	productStore types.ProductStore
	templates    types.Templates
}

func NewHandler(userStore types.UserStore, cartStore types.CartStore, productStore types.ProductStore, templates types.Templates) *Handler {
	return &Handler{userStore: userStore, cartStore: cartStore, productStore: productStore, templates: templates}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/cart", auth.NewHandler(h.userStore).WithJWT, h.GetCart)
	router.POST("/products/addtocart/:productID", auth.NewHandler(h.userStore).WithJWT, h.AddToCart)
}

func (h *Handler) GetCart(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(types.User)
	log.Println("user: ", user)
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
func (h *Handler) AddToCart(c *gin.Context) {
	u, _ := c.Get("user")
	user := u.(types.User)
	log.Println("user: ", user)
	var cart types.Cart
	quantity, err := strconv.Atoi(h.templates.GetDataFromForm(c, "quantity"))
	if err != nil {
		log.Println(err)
		return
	}
	id := c.Param("productID")
	log.Println("quantity", quantity)
	log.Println("id :", id)
	cart.Quantity = quantity
	cart.UserId = user.ID
	cart.Product_id = id
	h.cartStore.AddToCart(cart)
	if strings.HasPrefix(id, "a") {
		str := fmt.Sprintf("/products/accessory/%s", id)
		c.Writer.Header().Add("HX-Redirect", str)
	} else {
		str := fmt.Sprintf("/products/bicycle/%s", id)
		c.Writer.Header().Add("HX-Redirect", str)
	}
}
