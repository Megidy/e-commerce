package order

import (
	"log"

	order "github.com/Megidy/e-commerce/frontend/templates/order"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/Megidy/e-commerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	templates    types.Templates
	userStore    types.UserStore
	orderStore   types.OrderStore
	productStore types.ProductStore
	cartStore    types.CartStore
}

func NewHandler(templates types.Templates, userStore types.UserStore, orderStore types.OrderStore, productStore types.ProductStore, cartStore types.CartStore) *Handler {
	return &Handler{templates: templates, userStore: userStore, orderStore: orderStore, productStore: productStore, cartStore: cartStore}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authHandler *auth.Handler) {
	router.GET("/orders", authHandler.WithJWT, h.GetOrders)
	router.GET("/orders/:orderID", authHandler.WithJWT, h.GetOrderById)
	router.DELETE("/orders/:orderID/cancel", authHandler.WithJWT, h.CancelOrder)
	router.GET("/orders/confirm", authHandler.WithJWT, h.LoadConfirmOrderPage)
	router.POST("/orders/confirm/accept", authHandler.WithJWT, h.ConfirmOrder)
	router.POST("/orders/confirm/redirect", authHandler.WithJWT, h.Redirect)
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
	orderDetails, err := h.orderStore.GetOrderDetails(orderID, o.User_id)
	if err != nil {
		log.Println(err)
		return
	}
	order.OrderPage(o, ordersProducts, orderDetails, bicycles, accessories).Render(c.Request.Context(), c.Writer)
}

func (h *Handler) CancelOrder(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		return
	}
	orderID := c.Param("orderID")
	user := u.(types.User)
	err := h.orderStore.CancelOrder(orderID, user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	redirectURL := "/orders"
	c.Writer.Header().Add("HX-Redirect", redirectURL)
}

func (h *Handler) LoadConfirmOrderPage(c *gin.Context) {
	//show overview of order that is taken from cart
	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)
	log.Println("user: ", user)
	cart, err := h.cartStore.GetCart(user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	accessories, bicycles, totalPrice, err := h.productStore.GetAllProductsForCart(cart)
	if err != nil {
		log.Println(err)
		return
	}
	for _, accessory := range accessories {
		if accessory.Quantity > accessory.Accessory.Quantity {
			amount := accessory.Quantity - accessory.Accessory.Quantity
			err = h.cartStore.ChangeCartsProductQuantity(user.ID, accessory.Accessory.Id, "dec", amount)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	for _, bicycle := range bicycles {
		if bicycle.Quantity > bicycle.Bicycle.Quantity {
			amount := bicycle.Quantity - bicycle.Bicycle.Quantity
			err = h.cartStore.ChangeCartsProductQuantity(user.ID, bicycle.Bicycle.Id, "dec", amount)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
	//button to make order
	order.ConfirmOrderPage(accessories, bicycles, totalPrice).Render(c.Request.Context(), c.Writer)

}
func (h *Handler) ConfirmOrder(c *gin.Context) {
	//request data about user : name , email , phone number , address: country,city,street,house
	var payload types.OrderDetails
	var OrderDetails types.OrderDetails

	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)

	payload.FirstName = h.templates.GetDataFromForm(c, "name")
	payload.LastName = h.templates.GetDataFromForm(c, "lastname")
	payload.Email = h.templates.GetDataFromForm(c, "email")
	payload.PhoneNumber = h.templates.GetDataFromForm(c, "phonenumber")
	payload.Country = h.templates.GetDataFromForm(c, "country")
	payload.City = h.templates.GetDataFromForm(c, "city")
	payload.Street = h.templates.GetDataFromForm(c, "street")
	payload.House = h.templates.GetDataFromForm(c, "house")

	OrderDetails = types.OrderDetails{
		UserId:      user.ID,
		Order_Id:    uuid.NewString(),
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Country:     payload.Country,
		City:        payload.City,
		Street:      payload.Street,
		House:       payload.House,
	}
	err := h.orderStore.AddOrdersDetails(OrderDetails)
	if err != nil {
		log.Println(err)
	}
	RedirectURL := "/cart"
	var newOrder types.Order
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
		Order_id:   OrderDetails.Order_Id,
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
		err = h.productStore.ChangeProductsQuantity(ca.Product_id, "desc", ca.Quantity)
		if err != nil {
			log.Println(err)
			return
		}
	}

	c.Writer.Header().Add("HX-Redirect", RedirectURL)

}

func (h *Handler) Redirect(c *gin.Context) {
	c.Writer.Header().Add("HX-redirect", "/orders/confirm")
}
