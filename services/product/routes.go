package product

import (
	"fmt"
	"log"
	"strconv"

	templates "github.com/Megidy/e-commerce/frontend/templates/product"
	users "github.com/Megidy/e-commerce/frontend/templates/user"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/Megidy/e-commerce/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	templates    types.Templates
	userStore    types.UserStore
	productStore types.ProductStore
}

func NewHandler(templates types.Templates, userStore types.UserStore, productStore types.ProductStore) *Handler {
	return &Handler{
		templates:    templates,
		userStore:    userStore,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, authHandler *auth.Handler, managerHandler *auth.Manager) {
	//getting products
	router.GET("/products/accessories", h.GetAllAccessories)
	router.GET("/products/bicycles/", h.GetAllBicycles)
	router.GET("/products/bicycle/:bicycleID", h.GetBicycleById)
	router.GET("/products/accessory/:accessoryID", h.GetAccessoryById)
	//modifying products
	router.GET("/products/accessory/:accessoryID/modify", authHandler.WithJWT, managerHandler.WithManagerRole, h.LoadModifyAccessoryPage)
	router.GET("/products/bicycle/:bicycleID/modify", authHandler.WithJWT, managerHandler.WithManagerRole, h.LoadModifyBicyclePage)
	router.POST("/products/accessory/:accessoryID/modify/confirm", authHandler.WithJWT, managerHandler.WithManagerRole, h.ModifyAccessoryPageConfirm)
	router.POST("/products/bicycle/:bicycleID/modify/confirm", authHandler.WithJWT, managerHandler.WithManagerRole, h.ModifyBicyclePageConfirm)
	//redirecting products
	router.POST("/products/action/redirect", authHandler.WithJWT, managerHandler.WithManagerRole, h.ActionRedirector)
	//deleting products
	router.GET("/products/action/delete/:productID", authHandler.WithJWT, managerHandler.WithManagerRole, h.LoadDeleteConfirmation)
	router.DELETE("/products/action/delete/:productID/confirm", authHandler.WithJWT, managerHandler.WithManagerRole, h.DeleteProduct)
	//creating products
	router.GET("/products/accessory/action/add", authHandler.WithJWT, managerHandler.WithManagerRole, h.LoadAddAccessoryPage)
	router.GET("/products/bicycle/action/add", authHandler.WithJWT, managerHandler.WithManagerRole, h.LoadAddbicyclePage)
	router.POST("/products/accessory/action/add/confirm", authHandler.WithJWT, managerHandler.WithManagerRole, h.ConfirmAddingAccessory)
	router.POST("/products/bicycle/action/add/confirm", authHandler.WithJWT, managerHandler.WithManagerRole, h.ConfirmAddingBicycle)
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

func (h *Handler) GetBicycleById(c *gin.Context) {
	isAddingToCart := c.Request.URL.Query().Get("isAddingCar") == "true"
	id := c.Param("bicycleID")
	bicycle, err := h.productStore.GetBicycleById(id)
	if err != nil {
		log.Println(err)
		return
	}
	templates.LoadSingleBicycle(bicycle, isAddingToCart, "").Render(c.Request.Context(), c.Writer)
}
func (h *Handler) GetAccessoryById(c *gin.Context) {
	isAddingToCar := c.Request.URL.Query().Get("isAddingCar") == "true"
	id := c.Param("accessoryID")
	accessory, err := h.productStore.GetAccessoryById(id)
	if err != nil {
		log.Println(err)
		return
	}
	templates.LoadSingleAccessory(accessory, isAddingToCar, "").Render(c.Request.Context(), c.Writer)
}

func (h *Handler) ActionRedirector(c *gin.Context) {
	product := h.templates.GetDataFromForm(c, "product")
	action := h.templates.GetDataFromForm(c, "action")
	addproduct := h.templates.GetDataFromForm(c, "addproduct")
	switch action {
	case "modify":
		if utils.IsAccessory(product) {
			c.Writer.Header().Add("HX-Redirect", fmt.Sprintf("/products/accessory/%s/modify", product))
		} else {
			c.Writer.Header().Add("HX-Redirect", fmt.Sprintf("/products/bicycle/%s/modify", product))
		}

	case "delete":
		c.Writer.Header().Add("HX-Redirect", fmt.Sprintf("/products/action/delete/%s", product))
	}
	log.Println("product: ", product, " , action : ", action, ", addProduct : ", addproduct)
	switch addproduct {
	case "AddAccessory":
		c.Writer.Header().Add("HX-Redirect", "/products/accessory/action/add")
	case "AddBicycle":
		c.Writer.Header().Add("HX-Redirect", "/products/bicycle/action/add")
	}

}
func (h *Handler) DeleteProduct(c *gin.Context) {
	productID := c.Param("productID")
	err := h.productStore.DeleteProduct(productID)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("deleted product : ", productID)
	c.Writer.Header().Add("HX-Redirect", "/user/manager")
}
func (h *Handler) LoadDeleteConfirmation(c *gin.Context) {
	productID := c.Param("productID")
	users.LoadDeleteConfirmationPage(productID).Render(c.Request.Context(), c.Writer)
}

func (h *Handler) LoadAddAccessoryPage(c *gin.Context) {
	templates.LoadAddAccessoryPage("").Render(c.Request.Context(), c.Writer)
}

func (h *Handler) ConfirmAddingAccessory(c *gin.Context) {
	var Accessory types.Accessory
	id := h.templates.GetDataFromForm(c, "id")
	if !utils.IsAccessory(id) {
		templates.LoadAddAccessoryPage("id has to start with prefix 'a' ").Render(c.Request.Context(), c.Writer)
		return
	}
	ok, err := h.productStore.AccessoryAlreadyExists(id)
	if err != nil {
		templates.LoadAddAccessoryPage(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	if ok {
		log.Println("accessory already exists")
		templates.LoadAddAccessoryPage("accessory with this id already exists").Render(c.Request.Context(), c.Writer)
		return
	}

	name := h.templates.GetDataFromForm(c, "name")
	description := h.templates.GetDataFromForm(c, "description")
	q := h.templates.GetDataFromForm(c, "quantity")
	p := h.templates.GetDataFromForm(c, "price")
	category := h.templates.GetDataFromForm(c, "category")
	image := h.templates.GetDataFromForm(c, "image")
	price, err := strconv.ParseFloat(p, 32)
	if err != nil {
		templates.LoadAddAccessoryPage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	quantity, err := strconv.Atoi(q)
	if err != nil {
		templates.LoadAddAccessoryPage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	Accessory = types.Accessory{
		Id:          id,
		Name:        name,
		Description: description,
		Quantity:    quantity,
		Price:       float32(price),
		Category:    category,
		Image:       image,
	}
	err = h.productStore.AddAccessory(Accessory)
	if err != nil {
		templates.LoadAddAccessoryPage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	c.Writer.Header().Add("HX-Redirect", "/user/manager")

}
func (h *Handler) LoadAddbicyclePage(c *gin.Context) {
	templates.LoadAddBicyclePage("").Render(c.Request.Context(), c.Writer)
}
func (h *Handler) ConfirmAddingBicycle(c *gin.Context) {
	var Bicycle types.Bicycle
	id := h.templates.GetDataFromForm(c, "id")
	if !utils.IsBicycle(id) {
		templates.LoadAddBicyclePage("id has to start with prefix 'b'").Render(c.Request.Context(), c.Writer)
		return
	}
	ok, err := h.productStore.BicycleAlreadyExists(id)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	if ok {
		log.Println("accessory already exists")
		templates.LoadAddBicyclePage("bicycle with this id already exists").Render(c.Request.Context(), c.Writer)
		return
	}
	name := h.templates.GetDataFromForm(c, "name")
	model := h.templates.GetDataFromForm(c, "model")
	description := h.templates.GetDataFromForm(c, "description")
	typ := h.templates.GetDataFromForm(c, "type")
	size := h.templates.GetDataFromForm(c, "size")
	material := h.templates.GetDataFromForm(c, "material")
	q := h.templates.GetDataFromForm(c, "quantity")
	p := h.templates.GetDataFromForm(c, "price")
	image := h.templates.GetDataFromForm(c, "image")
	color := h.templates.GetDataFromForm(c, "color")
	w := h.templates.GetDataFromForm(c, "weight")
	r := h.templates.GetDataFromForm(c, "releaseyear")
	brakesystem := h.templates.GetDataFromForm(c, "brakesystem")
	g := h.templates.GetDataFromForm(c, "gears")
	brand := h.templates.GetDataFromForm(c, "brand")
	suspension := h.templates.GetDataFromForm(c, "suspension")
	ws := h.templates.GetDataFromForm(c, "wheelsize")
	framesize := h.templates.GetDataFromForm(c, "framesize")
	quantity, err := strconv.Atoi(q)
	if err != nil {
		log.Println(err)
		return
	}
	price, err := strconv.ParseFloat(p, 32)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	weight, err := strconv.Atoi(w)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	releaseyear, err := strconv.Atoi(r)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	gears, err := strconv.Atoi(g)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	wheelsize, err := strconv.Atoi(ws)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	Bicycle = types.Bicycle{
		Id:          id,
		Name:        name,
		Model:       model,
		Description: description,
		Type:        typ,
		Size:        size,
		Material:    material,
		Quantity:    quantity,
		Price:       float32(price),
		Image:       image,
		Color:       color,
		Weight:      float32(weight),
		ReleaseYear: releaseyear,
		BrakeSystem: brakesystem,
		Gears:       gears,
		Brand:       brand,
		Suspension:  suspension,
		WheelSize:   wheelsize,
		FrameSize:   framesize,
	}
	err = h.productStore.AddBicycle(Bicycle)
	if err != nil {
		templates.LoadAddBicyclePage(err.Error()).Render(c.Request.Context(), c.Writer)
		log.Println(err)
		return
	}
	c.Writer.Header().Add("HX-Redirect", "/user/manager")
}

func (h *Handler) LoadModifyAccessoryPage(c *gin.Context) {
	id := c.Param("accessoryID")
	templates.LoadModifyAccessoryPage(id).Render(c.Request.Context(), c.Writer)
}
func (h *Handler) ModifyAccessoryPageConfirm(c *gin.Context) {

	accessoryID := c.Param("accessoryID")
	name := h.templates.GetDataFromForm(c, "name")
	if name != "" {
		err := h.productStore.UpdateAccessory("name", name, accessoryID)
		if err != nil {
			log.Println(err)
			return
		}
	}
	description := h.templates.GetDataFromForm(c, "description")
	if description != "" {
		err := h.productStore.UpdateAccessory("description", description, accessoryID)
		if err != nil {
			log.Println(err)
			return
		}
	}
	q := h.templates.GetDataFromForm(c, "quantity")
	if q != "" {
		quantity, err := strconv.Atoi(q)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.productStore.UpdateAccessory("quantity", quantity, accessoryID)
		if err != nil {
			log.Println(err)
			return
		}
	}
	p := h.templates.GetDataFromForm(c, "price")
	if p != "" {
		p, err := strconv.ParseFloat(p, 32)
		if err != nil {
			log.Println(err)
			return
		}
		price := float32(p)
		err = h.productStore.UpdateAccessory("price", price, accessoryID)
		if err != nil {
			log.Println(err)
			return
		}
	}
	category := h.templates.GetDataFromForm(c, "category")
	if category != "" {
		err := h.productStore.UpdateAccessory("category", category, accessoryID)
		if err != nil {
			log.Println(err)
			return
		}
	}
	image := h.templates.GetDataFromForm(c, "image")
	if image != "" {
		err := h.productStore.UpdateAccessory("image", image, accessoryID)
		if err != nil {
			log.Println(err)
			return
		}
	}
	c.Writer.Header().Add("HX-Redirect", "/user/manager")
}

func (h *Handler) LoadModifyBicyclePage(c *gin.Context) {
	id := c.Param("bicycleID")
	templates.LoadModifyBicyclePage(id).Render(c.Request.Context(), c.Writer)
}
func (h *Handler) ModifyBicyclePageConfirm(c *gin.Context) {

	id := c.Param("bicycleID")

	name := h.templates.GetDataFromForm(c, "name")
	if name != "" {
		err := h.productStore.UpdateBicycle("name", name, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	model := h.templates.GetDataFromForm(c, "model")
	if model != "" {
		err := h.productStore.UpdateBicycle("model", model, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	description := h.templates.GetDataFromForm(c, "description")
	if description != "" {
		err := h.productStore.UpdateBicycle("description", description, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	typ := h.templates.GetDataFromForm(c, "type")
	if typ != "" {
		err := h.productStore.UpdateBicycle("type", typ, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	size := h.templates.GetDataFromForm(c, "size")
	if size != "" {
		err := h.productStore.UpdateBicycle("size", size, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	material := h.templates.GetDataFromForm(c, "material")
	if material != "" {
		err := h.productStore.UpdateBicycle("material", material, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	q := h.templates.GetDataFromForm(c, "quantity")
	if q != "" {
		quantity, err := strconv.Atoi(q)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.productStore.UpdateBicycle("quantity", quantity, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	p := h.templates.GetDataFromForm(c, "price")
	if p != "" {
		p, err := strconv.ParseFloat(p, 32)
		if err != nil {
			log.Println(err)
			return
		}
		price := float32(p)
		err = h.productStore.UpdateBicycle("price", price, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	image := h.templates.GetDataFromForm(c, "image")
	if image != "" {
		err := h.productStore.UpdateBicycle("image", image, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	color := h.templates.GetDataFromForm(c, "color")
	if color != "" {
		err := h.productStore.UpdateBicycle("color", color, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	w := h.templates.GetDataFromForm(c, "weight")
	if w != "" {
		w, err := strconv.ParseFloat(w, 32)
		if err != nil {
			log.Println(err)
			return
		}
		weight := float32(w)
		err = h.productStore.UpdateBicycle("weight", weight, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	r := h.templates.GetDataFromForm(c, "releaseyear")
	if r != "" {
		releaseyear, err := strconv.Atoi(r)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.productStore.UpdateBicycle("releaseyear", releaseyear, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	brakesystem := h.templates.GetDataFromForm(c, "brakesystem")
	if brakesystem != "" {
		err := h.productStore.UpdateBicycle("brakesystem", brakesystem, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	g := h.templates.GetDataFromForm(c, "gears")
	if g != "" {
		gears, err := strconv.Atoi(g)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.productStore.UpdateBicycle("gears", gears, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	brand := h.templates.GetDataFromForm(c, "brand")
	if brand != "" {
		err := h.productStore.UpdateBicycle("brand", brand, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	suspension := h.templates.GetDataFromForm(c, "suspension")
	if suspension != "" {
		err := h.productStore.UpdateBicycle("suspension", suspension, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	ws := h.templates.GetDataFromForm(c, "wheelsize")
	if ws != "" {
		wheelSize, err := strconv.Atoi(ws)
		if err != nil {
			log.Println(err)
			return
		}
		err = h.productStore.UpdateBicycle("wheelsize", wheelSize, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	framesize := h.templates.GetDataFromForm(c, "framesize")
	if framesize != "" {
		err := h.productStore.UpdateBicycle("framesize", framesize, id)
		if err != nil {
			log.Println(err)
			return
		}
	}
	c.Writer.Header().Add("HX-Redirect", "/user/manager")

}
