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

	_, err := s.db.Exec("insert into orders(order_id,user_id,status) values(?,?,?)", order.Order_id, order.User_id, order.Status)
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
