package product

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
