package order

import (
	"log"

	"github.com/Megidy/e-commerce/frontend/templates/order"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/Megidy/e-commerce/utils"
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
	router.GET("/orders/:orderID", authHandler.WithJWT, h.GetOrderById)
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
	_, _, totalPrice, err := h.productStore.GetAllProductsForCart(cart)
	if err != nil {
		log.Println(err)
		return
	}
	newOrder = types.Order{
		Order_id:   uuid.NewString(),
		User_id:    user.ID,
		Status:     "Processing ",
		TotalPrice: totalPrice,
	}
	err = h.orderStore.CreateOrder(newOrder)
	if err != nil {
		log.Println(err)
		return
	}

	for _, ca := range cart {
		if ca.Quantity == 0 {

			c.Writer.Header().Add("HX-Redirect", RedirectURL)
			return
		}
		log.Println(ca.Product_id)
		if utils.IsAccessory(ca.Product_id) {
			acc, err := h.productStore.GetAccessoryById(ca.Product_id)
			if err != nil {
				log.Println(err)
				return
			}
			err = h.orderStore.CreateOrdersProducts(newOrder, ca, acc.Price*float32(ca.Quantity))
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("accessory: ", acc)
		} else {
			bic, err := h.productStore.GetBicycleById(ca.Product_id)
			if err != nil {
				log.Println(err)
				return
			}
			err = h.orderStore.CreateOrdersProducts(newOrder, ca, bic.Price*float32(ca.Quantity))
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(bic)
		}
		err = h.cartStore.DeleteFromCart(ca.UserId, ca.Product_id)
		if err != nil {
			log.Println(err)
			return
		}
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

func (h *Handler) GetOrderById(c *gin.Context) {
	var accessories []types.Accessory
	var bicycles []types.Bicycle
	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found ")
	}
	user := u.(types.User)
	//get order by id
	orderID := c.Param("orderID")
	//show that order
	o, err := h.orderStore.GetOrderById(user.ID, orderID)
	if err != nil {
		log.Println(err)
		return
	}
	ordersProducts, err := h.orderStore.GetOrdersProducts(user.ID, orderID)
	if err != nil {
		log.Println(err)
		return
	}

	for _, product := range ordersProducts {
		if utils.IsAccessory(product.Product_id) {
			acc, err := h.productStore.GetAccessoryById(product.Product_id)
			if err != nil {
				log.Println(err)
				return
			}
			accessories = append(accessories, acc)
		} else {
			bic, err := h.productStore.GetBicycleById(product.Product_id)
			if err != nil {
				log.Println(err)
				return
			}
			bicycles = append(bicycles, bic)
		}

	}
	order.OrderPage(o, ordersProducts, bicycles, accessories).Render(c.Request.Context(), c.Writer)

}
