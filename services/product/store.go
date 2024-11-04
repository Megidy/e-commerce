package product

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Megidy/e-commerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAllAccessories() ([]types.Accessory, error) {
	var accessories []types.Accessory
	rows, err := s.db.Query("select * from accessories")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var a types.Accessory
		err := rows.Scan(&a.Id, &a.Name, &a.Description, &a.Quantity, &a.Price, &a.Image)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		accessories = append(accessories, a)
	}
	log.Println(accessories)
	return accessories, nil
}

func (s *Store) GetAllBicycles() ([]types.Bicycle, error) {
	var bicycles []types.Bicycle
	rows, err := s.db.Query("select * from bicycles")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for rows.Next() {
		var b types.Bicycle
		err := rows.Scan(&b.Id, &b.Name, &b.Model, &b.Description, &b.Type, &b.Size, &b.Material, &b.Quantity, &b.Price, &b.Image)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		bicycles = append(bicycles, b)
	}
	log.Println(bicycles)
	return bicycles, nil
}

func (s *Store) GetBicycleById(id string) (types.Bicycle, error) {
	var b types.Bicycle
	row, err := s.db.Query("select * from bicycles where id =?", id)
	if err != nil {
		return types.Bicycle{}, err
	}
	for row.Next() {

		err := row.Scan(&b.Id, &b.Name, &b.Model, &b.Description, &b.Type, &b.Size, &b.Material, &b.Quantity, &b.Price, &b.Image)
		if err != nil {
			return types.Bicycle{}, err

		}
	}
	return b, nil
}
func (s *Store) GetAccessoryById(id string) (types.Accessory, error) {
	var a types.Accessory
	row, err := s.db.Query("select * from accessories where id =?", id)
	if err != nil {
		return types.Accessory{}, err
	}
	for row.Next() {

		err := row.Scan(&a.Id, &a.Name, &a.Description, &a.Quantity, &a.Price, &a.Image)
		if err != nil {
			return types.Accessory{}, err

		}
	}
	return a, nil
}

func (s *Store) GetAllProducts(carts []types.Cart) ([]types.Accessory, []types.Bicycle, error) {
	var accessories []types.Accessory
	var bicycle []types.Bicycle
	log.Println("carts : ", carts)
	for _, cart := range carts {
		log.Println("entered loop")
		var b types.Bicycle
		if strings.HasPrefix(cart.Product_id, "b") {
			log.Println("entered b")
			rows, err := s.db.Query("select * from bicycles where id = ?", cart.Product_id)
			if err != nil {
				return nil, nil, err
			}
			log.Println("queried")
			for rows.Next() {
				err := rows.Scan(&b.Id, &b.Name, &b.Model, &b.Description, &b.Type, &b.Size, &b.Material, &b.Quantity, &b.Price, &b.Image)
				if err != nil {
					return nil, nil, err
				}
				log.Println("scanned")
				bicycle = append(bicycle, b)
			}
		} else {
			var a types.Accessory
			rows, err := s.db.Query("select * from accessories where id = ?", cart.Product_id)
			if err != nil {
				return nil, nil, err
			}
			for rows.Next() {
				err := rows.Scan(&a.Id, &a.Name, &a.Description, &a.Quantity, &a.Price, &a.Image)
				if err != nil {
					return nil, nil, err
				}
				accessories = append(accessories, a)
			}
		}

	}

	return accessories, bicycle, nil
}
