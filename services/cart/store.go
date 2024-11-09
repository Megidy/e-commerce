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

func (s *Store) DeleteFromCart(userID, productID string) error {
	_, err := s.db.Exec("delete from cart where user_id=? and product_id=?", userID, productID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CheckIfProductInCart(userID, productID string) (bool, error) {
	row, err := s.db.Query("select * from cart where user_id=? and product_id=?", userID, productID)
	if err != nil {
		return false, err
	}
	for row.Next() {
		var c types.Cart
		err = row.Scan(&c.UserId, &c.Product_id, &c.Quantity, &c.Created)
		if err != nil {
			if err == sql.ErrNoRows {
				return false, nil
			}
			return false, err
		}
		if productID == c.Product_id {
			return true, nil
		}

	}
	return false, nil
}
func (s *Store) ChangeProductsQuantity(userID, ProductID string, amount int) error {
	_, err := s.db.Exec("update cart set quantity = quantity+? where user_id=? and product_id =?", amount, userID, ProductID)
	if err != nil {
		return err
	}
	return nil
}
