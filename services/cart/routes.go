package cart

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	carts "github.com/Megidy/e-commerce/frontend/templates/cart"
	"github.com/Megidy/e-commerce/utils"

	product "github.com/Megidy/e-commerce/frontend/templates/product"
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

	carts.LoadCart(bicycles, accessories).Render(c.Request.Context(), c.Writer)

}
func (h *Handler) AddToCart(c *gin.Context) {

	//setting up url query
	isAddingToCart := c.Request.URL.Query().Get("isAddingCar") == "true"
	var quantityOfproduct int
	//getting params
	productID := c.Param("productID")
	var cart types.Cart

	//getting user for cart
	u, _ := c.Get("user")
	user := u.(types.User)
	//converting to int
	//if payload is incorrect than ->bad request
	str := h.templates.GetDataFromForm(c, "quantity")
	quantity, err := strconv.Atoi(str)

	//checking if payload is correct for both
	if err != nil {
		if utils.IsAccessory(productID) {
			acc, err := h.productStore.GetAccessoryById(productID)
			if err != nil {
				log.Println(err)
				return
			}

			product.LoadSingleAccessory(acc, isAddingToCart, "Bad Request").Render(c.Request.Context(), c.Writer)
			log.Println(err)
			return
		} else {
			bic, err := h.productStore.GetBicycleById(productID)
			if err != nil {
				log.Println(err)
				return
			}
			product.LoadSingleBicycle(bic, isAddingToCart, "Bad Request").Render(c.Request.Context(), c.Writer)
			return
		}

	}

	//getting product from db for both
	if utils.IsAccessory(productID) {
		acc, err := h.productStore.GetAccessoryById(productID)
		if err != nil {
			log.Println(err)
			return
		}
		quantityOfproduct = acc.Quantity

	} else {
		bic, err := h.productStore.GetBicycleById(productID)
		log.Println("bicycle: ", bic)
		if err != nil {
			log.Println(err)
			return
		}
		quantityOfproduct = bic.Quantity
	}

	//checking if quantity of product > quantity of order for both
	if quantity > quantityOfproduct {
		if utils.IsAccessory(productID) {
			acc, err := h.productStore.GetAccessoryById(productID)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("quantity of order is bigger than quantity of product")
			product.LoadSingleAccessory(acc, isAddingToCart, "Bad Request").Render(c.Request.Context(), c.Writer)
			return
		} else {
			bic, err := h.productStore.GetBicycleById(productID)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("quantity of order is bigger than quantity of product")
			product.LoadSingleBicycle(bic, isAddingToCart, "Bad Request").Render(c.Request.Context(), c.Writer)
			return
		}

	}

	cart.Quantity = quantity
	cart.UserId = user.ID
	cart.Product_id = productID

	err = h.cartStore.AddToCart(cart)
	if err != nil {
		log.Println(err)
		return
	}
	if utils.IsAccessory(productID) {
		str := fmt.Sprintf("/products/accessory/%s", productID)
		c.Writer.Header().Add("HX-Redirect", str)
	} else {
		str := fmt.Sprintf("/products/bicycle/%s", productID)
		c.Writer.Header().Add("HX-Redirect", str)

	}

}
