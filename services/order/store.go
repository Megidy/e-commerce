package order

import (
	"database/sql"
	"log"

	"github.com/Megidy/e-commerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) error {

	_, err := s.db.Exec("insert into orders(order_id,user_id,status,total_price) values(?,?,?,?)", order.Order_id, order.User_id, order.Status, order.TotalPrice)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateOrdersProducts(order types.Order, cart types.Cart, totalPrice float32) error {
	_, err := s.db.Exec("insert into orders_products(user_id,order_id,product_id,quantity,total_price) values(?,?,?,?,?)", cart.UserId, order.Order_id, cart.Product_id, cart.Quantity, totalPrice)

	if err != nil {
		return err
	}
	log.Printf("Inserting order_id: %s, product_id: %s, quantity: %d, total_price: %.2f\n", order.Order_id, cart.Product_id, cart.Quantity, totalPrice)

	return nil
}

func (s *Store) GetOrders(userID string) ([]types.Order, error) {
	var orders []types.Order
	rows, err := s.db.Query("select * from orders where user_id=?", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var o types.Order
		err = rows.Scan(&o.Order_id, &o.User_id, &o.Status, &o.Created, &o.TotalPrice)
		if err != nil {
			return nil, err

		}
		orders = append(orders, o)
	}
	return orders, nil

}
func (s *Store) GetOrderById(userID, orderId string) (types.Order, error) {
	row, err := s.db.Query("select * from orders where order_id=? and user_id=?", orderId, userID)
	if err != nil {
		return types.Order{}, nil
	}
	var order types.Order

	for row.Next() {
		err = row.Scan(&order.Order_id, &order.User_id, &order.Status, &order.Created, &order.TotalPrice)
		if err != nil {
			return types.Order{}, nil
		}
	}
	return order, nil
}
func (s *Store) GetOrdersProducts(userID, orderID string) ([]types.OrderProduct, error) {
	var products []types.OrderProduct
	rows, err := s.db.Query("select * from orders_products where user_id=? and order_id=?", userID, orderID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product types.OrderProduct
		err = rows.Scan(&product.User_id, &product.Order_id, &product.Product_id, &product.Quantity, &product.TotalPrice)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) CancelOrder(orderID, userID string) error {
	status := "canceled"
	_, err := s.db.Exec("update orders set status=? where order_id =? and user_id=?", status, orderID, userID)
	if err != nil {
		return err
	}
	return nil
}
