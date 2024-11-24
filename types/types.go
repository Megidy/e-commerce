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
	DeleteProduct(productID string) error
	AddAccessory(accessory Accessory) error
	AddBicycle(bicycle Bicycle) error
	AccessoryAlreadyExists(id string) (bool, error)
	BicycleAlreadyExists(id string) (bool, error)
}

type CartStore interface {
	GetCart(userID string) ([]Cart, error)
	AddToCart(cart Cart) error
	DeleteFromCart(userID, productID string) error
	ProductInCart(userID, productID string) (bool, error)
	ChangeCartsProductQuantity(userID, ProductID, action string, amount int) error
}
type OrderStore interface {
	CreateOrder(order Order) error
	GetOrders(userID string) ([]Order, error)
	GetOrderById(userID, orderId string) (Order, error)
	CreateOrdersProducts(order Order, cart Cart, totalPrice float32) error
	GetOrdersProducts(userID, orderID string) ([]OrderProduct, error)
	CancelOrder(orderID, userID string) error
	AddOrdersDetails(orderDetails OrderDetails) error
	GetOrderDetails(orderID, userID string) (OrderDetails, error)
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
	User_id    string
	Order_id   string
	Product_id string
	Quantity   int
	TotalPrice float32
}

type OrderDetails struct {
	UserId      string
	Order_Id    string
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Country     string
	City        string
	Street      string
	House       string
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
