package types

import "github.com/gin-gonic/gin"

type User struct {
	ID       string
	Name     string
	LastName string
	Email    string
	Password string
	Created  string
	Role     string
}

type CreateUserPayload struct {
	Name     string
	LastName string
	Email    string
	Password string
}

type LogInPayload struct {
	Email    string
	Password string
}

type UserStore interface {
	CreateUser(user *User) error
	AlreadyExists(user *User) (bool, error)
	GetUserByEmail(email string) (User, error)
	GetUserById(id string) (User, error)
}

type Templates interface {
	LoadResponse(c *gin.Context, statusCode int, file string, response map[string]any)
	ExecuteTemplate(c *gin.Context, filePath string, data map[string]any) error
	ExecuteSpecificTemplate(c *gin.Context, name, filePath string, data map[string]any) error
	GetDataFromForm(c *gin.Context, key string) string
}
type ProductStore interface {
	GetAllAccessories() ([]Accessory, error)
	GetAllBicycles() ([]Bicycle, error)
	GetBicycleById(id string) (Bicycle, error)
	GetAccessoryById(id string) (Accessory, error)
	GetAllProducts(carts []Cart) ([]CartAccessoryResponse, []CartBicycleResponse, float32, error)
	ChangeProductsQuantity(productID, action string, amount int) error
}

type CartStore interface {
	GetCart(userID string) ([]Cart, error)
	AddToCart(cart Cart) error
	DeleteFromCart(userID, productID string) error
	CheckIfProductInCart(userID, productID string) (bool, error)
	ChangeProductsQuantity(userID, ProductID string, amount int) error
}
type Accessory struct {
	Id          string
	Name        string
	Description string
	Quantity    int
	Price       float32
	Category    string
	Image       string
}

type Bicycle struct {
	Id          string
	Name        string
	Model       string
	Description string
	Type        string
	Size        string
	Material    string
	Quantity    int
	Price       float32
	Image       string
}

type Cart struct {
	UserId     string
	Product_id string
	Quantity   int
	Created    string
}

type CartBicycleResponse struct {
	Bicycle        *Bicycle
	Quantity       int
	PriceOfBicycle float32
}

type CartAccessoryResponse struct {
	Accessory        *Accessory
	Quantity         int
	PriceOfAccessory float32
}
