package order

import (
	"log"

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
	return &Handler{userStore: userStore, orderStore: orderStore, cartStore: cartStore}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/orders/createorder", auth.NewHandler(h.userStore).WithJWT, h.CreateOrder)
	router.GET("/orders")
	router.GET("/orders/:orderID")
	router.GET("orders/:orderID/:productID")
	router.DELETE("/orders/:orderID/cancel")

}

func (h *Handler) CreateOrder(c *gin.Context) {
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
	_, _, totalPrice, err := h.productStore.GetAllProductsForCart(cart)
	if err != nil {
		log.Println(err)
		return
	}
	for _, c := range cart {
		err = h.cartStore.DeleteFromCart(c.UserId, c.Product_id)
		if err != nil {
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
	str := "/cart"
	c.Writer.Header().Add("HX-Redirect", str)

}
