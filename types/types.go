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
	GetDataFromForm(c *gin.Context, key string) string
}
type ProductStore interface {
	GetAllAccessories() ([]Accessory, error)
	GetAllBicycles() ([]Bicycle, error)
	GetBicycleById(id string) (Bicycle, error)
	GetAccessoryById(id string) (Accessory, error)
	GetAllProductsForCart(carts []Cart) ([]CartAccessoryResponse, []CartBicycleResponse, float32, error)
	ChangeProductsQuantity(productID, action string, amount int) error
}

type CartStore interface {
	GetCart(userID string) ([]Cart, error)
	AddToCart(cart Cart) error
	DeleteFromCart(userID, productID string) error
	ProductInCart(userID, productID string) (bool, error)
	ChangeCartsProductQuantity(userID, ProductID, action string, amount int) error
}
type OrderStore interface {
	CreateOrder(order Order, cart []Cart) error
	GetOrders(userID string) ([]Order, error)
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
	Color       string
	Weight      float32
	ReleaseYear int
	BrakeSystem string
	Gears       int
	Brand       string
	Suspension  string
	WheelSize   int
	FrameSize   string
}

type Cart struct {
	UserId     string
	Product_id string
	Quantity   int
	Created    string
}

type Order struct {
	Order_id   string
	User_id    string
	Status     string
	Created    string
	TotalPrice float32
}

type OrderProduct struct {
	Order_id   string
	Product_id string
	Quantity   int
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
