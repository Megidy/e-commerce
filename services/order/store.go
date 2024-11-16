package order

import (
	"database/sql"

	"github.com/Megidy/e-commerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order, cart []types.Cart) error {

	_, err := s.db.Exec("insert into orders(order_id,user_id,status,total_price) values(?,?,?,?)", order.Order_id, order.User_id, order.Status, order.TotalPrice)
	if err != nil {
		return err
	}

	for _, c := range cart {
		_, err := s.db.Exec("insert into orders_products(order_id,product_id,quantity) values(?,?,?)", order.Order_id, c.Product_id, c.Quantity)
		if err != nil {
			return err

		}
	}
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
