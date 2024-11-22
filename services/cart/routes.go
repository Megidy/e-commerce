package cart

import (
	"fmt"
	"log"
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

func (h *Handler) RegisterRoutes(router gin.IRouter, authHandler *auth.Handler) {
	router.GET("/cart", authHandler.WithJWT, h.GetCart)
	router.POST("/products/addtocart/:productID", authHandler.WithJWT, h.AddToCart)
	router.DELETE("/cart/deletefromcart/:productID", authHandler.WithJWT, h.DeleteFromCart)
}

func (h *Handler) GetCart(c *gin.Context) {
	u, _ := c.Get("user")
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
				return
			}
		}
	}
	for _, bicycle := range bicycles {
		if bicycle.Quantity > bicycle.Bicycle.Quantity {
			amount := bicycle.Quantity - bicycle.Bicycle.Quantity
			err = h.cartStore.ChangeCartsProductQuantity(user.ID, bicycle.Bicycle.Id, "dec", amount)
			if err != nil {
				return
			}
		}
	}

	carts.LoadCart(accessories, bicycles, totalPrice).Render(c.Request.Context(), c.Writer)

}
func (h *Handler) AddToCart(c *gin.Context) {

	//setting up url query
	isAddingToCart := c.Request.URL.Query().Get("isAddingCar") == "true"
	var quantityOfproduct int
	var cart types.Cart
	//getting params
	productID := c.Param("productID")

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

	//Checking if quantity of payload  < quantity in cart
	ca, err := h.cartStore.GetCart(user.ID)
	if err != nil {
		log.Println(err)
		return
	}
	for _, car := range ca {
		if car.Product_id == productID {
			if utils.IsAccessory(productID) {
				acc, err := h.productStore.GetAccessoryById(productID)
				if err != nil {
					log.Println(err)
					return
				}
				if quantity+car.Quantity > acc.Quantity {
					product.LoadSingleAccessory(acc, isAddingToCart, "cart is overload for this product").Render(c.Request.Context(), c.Writer)
					log.Println(err)
					return
				}

			} else {
				bic, err := h.productStore.GetBicycleById(productID)
				if err != nil {
					log.Println(err)
					return
				}
				if quantity+car.Quantity > bic.Quantity {
					product.LoadSingleBicycle(bic, isAddingToCart, "cart is overload for this product").Render(c.Request.Context(), c.Writer)
					log.Println(err)
					return
				}
			}
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
			product.LoadSingleAccessory(acc, isAddingToCart, "Big Order").Render(c.Request.Context(), c.Writer)
			return
		} else {
			bic, err := h.productStore.GetBicycleById(productID)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("quantity of order is bigger than quantity of product")
			product.LoadSingleBicycle(bic, isAddingToCart, "Big Order").Render(c.Request.Context(), c.Writer)
			return
		}

	}

	cart.Quantity = quantity
	cart.UserId = user.ID
	cart.Product_id = productID
	ok, err := h.cartStore.ProductInCart(user.ID, productID)

	if err != nil {
		product.LoadSingleBicycle(types.Bicycle{}, isAddingToCart, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	if ok {
		err = h.cartStore.ChangeCartsProductQuantity(cart.UserId, productID, "inc", quantity)
		if err != nil {
			product.LoadSingleBicycle(types.Bicycle{}, isAddingToCart, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
	} else {
		err = h.cartStore.AddToCart(cart)
		if err != nil {
			product.LoadSingleBicycle(types.Bicycle{}, isAddingToCart, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
	}

	if utils.IsAccessory(productID) {
		str := fmt.Sprintf("/products/accessory/%s", productID)
		c.Writer.Header().Add("HX-Redirect", str)
	} else {
		str := fmt.Sprintf("/products/bicycle/%s", productID)
		c.Writer.Header().Add("HX-Redirect", str)

	}

}
func (h *Handler) DeleteFromCart(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)
	productID := c.Param("productID")
	err := h.cartStore.DeleteFromCart(user.ID, productID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("deleted from cart accessory: ", productID)

}
