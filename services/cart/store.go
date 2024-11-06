package cart

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

func (s *Store) GetCart(userID string) ([]types.Cart, error) {
	var cart []types.Cart
	rows, err := s.db.Query("select * from cart where user_id=?", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var c types.Cart

		err := rows.Scan(&c.UserId, &c.Product_id, &c.Quantity, &c.Created)
		if err != nil {
			return nil, err
		}
		cart = append(cart, c)
	}
	return cart, nil

}

func (s *Store) AddToCart(cart types.Cart) error {
	_, err := s.db.Exec("insert into cart(user_id,product_id,quantity) values (?,?,?)", cart.UserId, cart.Product_id, cart.Quantity)
	if err != nil {
		return err
	}
	return nil
}
