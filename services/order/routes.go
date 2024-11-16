package order

import (
	"log"

	"github.com/Megidy/e-commerce/frontend/templates/order"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userStore    types.UserStore
	orderStore   types.OrderStore
	productStore types.ProductStore
	cartStore    types.CartStore
}

func NewHandler(userStore types.UserStore, orderStore types.OrderStore, productStore types.ProductStore, cartStore types.CartStore) *Handler {
	return &Handler{userStore: userStore, orderStore: orderStore, productStore: productStore, cartStore: cartStore}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authHandler *auth.Handler) {
	router.POST("/orders/createorder", authHandler.WithJWT, h.CreateOrder)
	router.GET("/orders", authHandler.WithJWT, h.GetOrders)
	router.GET("/orders/:orderID")
	router.GET("orders/:orderID/:productID")
	router.DELETE("/orders/:orderID/cancel")

}

func (h *Handler) CreateOrder(c *gin.Context) {
	RedirectURL := "/cart"
	var newOrder types.Order
	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)
	cart, err := h.cartStore.GetCart(user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if cart == nil {
		c.Writer.Header().Add("HX-Redirect", RedirectURL)
		return
	}
	log.Println("cart : ", cart)
	log.Println("user : ", user)
	acc, bic, totalPrice, err := h.productStore.GetAllProductsForCart(cart)
	log.Println("acc: ", acc, "bic: ", bic)
	if err != nil {
		log.Println(err)
		return
	}
	for _, ca := range cart {
		if ca.Quantity == 0 {

			c.Writer.Header().Add("HX-Redirect", RedirectURL)
			return
		}
		err = h.cartStore.DeleteFromCart(ca.UserId, ca.Product_id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	newOrder = types.Order{
		Order_id:   uuid.NewString(),
		User_id:    user.ID,
		Status:     "Processing ",
		TotalPrice: totalPrice,
	}
	err = h.orderStore.CreateOrder(newOrder, cart)
	if err != nil {
		log.Println(err)
		return
	}
	c.Writer.Header().Add("HX-Redirect", RedirectURL)

}

func (h *Handler) GetOrders(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		log.Println("user was not found ")
		return
	}
	user := u.(types.User)
	orders, err := h.orderStore.GetOrders(user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	order.LoadOrders(orders).Render(c.Request.Context(), c.Writer)
}
