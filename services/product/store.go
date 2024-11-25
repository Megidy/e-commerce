package product

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Megidy/e-commerce/types"
	"github.com/Megidy/e-commerce/utils"
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
		err := rows.Scan(&a.Id, &a.Name, &a.Description, &a.Quantity, &a.Price, &a.Category, &a.Image)
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
		err := rows.Scan(&b.Id, &b.Name, &b.Model, &b.Description, &b.Type, &b.Size, &b.Material, &b.Quantity, &b.Price, &b.Image, &b.Color, &b.Weight, &b.ReleaseYear, &b.BrakeSystem, &b.Gears, &b.Brand, &b.Suspension, &b.WheelSize, &b.FrameSize)
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

		err := row.Scan(&b.Id, &b.Name, &b.Model, &b.Description, &b.Type, &b.Size, &b.Material, &b.Quantity, &b.Price, &b.Image, &b.Color, &b.Weight, &b.ReleaseYear, &b.BrakeSystem, &b.Gears, &b.Brand, &b.Suspension, &b.WheelSize, &b.FrameSize)
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

		err := row.Scan(&a.Id, &a.Name, &a.Description, &a.Quantity, &a.Price, &a.Category, &a.Image)
		if err != nil {
			return types.Accessory{}, err

		}
	}
	return a, nil
}

func (s *Store) GetAllProductsForCart(carts []types.Cart) ([]types.CartAccessoryResponse, []types.CartBicycleResponse, float32, error) {
	var accessories []types.CartAccessoryResponse
	var bicycles []types.CartBicycleResponse
	var totalPrice float32
	log.Println("carts : ", carts)
	for _, cart := range carts {
		log.Println("entered loop")
		if utils.IsAccessory(cart.Product_id) {
			var accessory types.CartAccessoryResponse

			var PriceOfAllAccessories float32
			row, err := s.db.Query("select * from accessories where id=?", cart.Product_id)
			if err != nil {
				return nil, nil, 0, err
			}
			for row.Next() {
				var a types.Accessory
				err := row.Scan(&a.Id, &a.Name, &a.Description, &a.Quantity, &a.Price, &a.Category, &a.Image)
				if err != nil {
					return nil, nil, 0, err
				}
				PriceOfAllAccessories = a.Price * float32(cart.Quantity)
				totalPrice = totalPrice + PriceOfAllAccessories
				accessory.Accessory = &a
				accessory.Quantity = cart.Quantity
				accessory.PriceOfAccessory = PriceOfAllAccessories
				accessories = append(accessories, accessory)
			}
			log.Println("accessories : ", accessories, ". Accessories PriceOfAccessory : ", accessory.PriceOfAccessory)
		} else {
			var bicycle types.CartBicycleResponse
			var PriceOfAllBicycles float32
			row, err := s.db.Query("select * from bicycles where id =? ", cart.Product_id)
			if err != nil {
				return nil, nil, 0, err
			}
			for row.Next() {
				var b types.Bicycle
				err := row.Scan(&b.Id, &b.Name, &b.Model, &b.Description, &b.Type, &b.Size, &b.Material, &b.Quantity, &b.Price, &b.Image, &b.Color, &b.Weight, &b.ReleaseYear, &b.BrakeSystem, &b.Gears, &b.Brand, &b.Suspension, &b.WheelSize, &b.FrameSize)
				if err != nil {
					return nil, nil, 0, err
				}
				PriceOfAllBicycles = b.Price * float32(cart.Quantity)
				totalPrice = totalPrice + PriceOfAllBicycles
				bicycle.Bicycle = &b
				bicycle.Quantity = cart.Quantity
				bicycle.PriceOfBicycle = PriceOfAllBicycles
				bicycles = append(bicycles, bicycle)
			}
		}

	}

	return accessories, bicycles, totalPrice, nil
}
func (s *Store) ChangeProductsQuantity(productID, action string, amount int) error {
	if action == "inc" {
		if utils.IsAccessory(productID) {
			_, err := s.db.Exec("update accessories set quantity = quantity+? where id =?", amount, productID)
			if err != nil {
				return err
			}
		} else {
			_, err := s.db.Exec("update bicycles set quantity = quantity+? where id =?", amount, productID)
			if err != nil {
				return err
			}
		}

	} else {
		if utils.IsAccessory(productID) {
			_, err := s.db.Exec("update accessories set quantity = quantity-? where id =?", amount, productID)
			if err != nil {
				return err
			}
		} else {
			_, err := s.db.Exec("update bicycles set quantity = quantity-? where id =?", amount, productID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Store) DeleteProduct(productID string) error {
	var err error
	if utils.IsAccessory(productID) {
		_, err = s.db.Exec("delete from accessories where id=?", productID)
	} else {
		_, err = s.db.Exec("delete from bicycles where id=?", productID)
	}
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) AddAccessory(accessory types.Accessory) error {
	_, err := s.db.Exec("insert into accessories values(?,?,?,?,?,?,?)", accessory.Id, accessory.Name, accessory.Description, accessory.Quantity, accessory.Price, accessory.Category, accessory.Image)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) AddBicycle(bicycle types.Bicycle) error {
	_, err := s.db.Exec("insert into bicycles values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", bicycle.Id, bicycle.Name, bicycle.Model, bicycle.Description, bicycle.Type, bicycle.Size, bicycle.Material, bicycle.Quantity, bicycle.Price, bicycle.Image, bicycle.Color, bicycle.Weight, bicycle.ReleaseYear, bicycle.BrakeSystem, bicycle.Gears, bicycle.Brand, bicycle.Suspension, bicycle.WheelSize, bicycle.FrameSize)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) AccessoryAlreadyExists(id string) (bool, error) {
	var accID string
	row, err := s.db.Query("select id from accessories where id=?", id)
	if err != nil {
		return true, err
	}
	for row.Next() {

		err = row.Scan(&accID)
		if err != nil {
			return true, err
		}
	}
	if accID != "" || accID == id {
		return true, nil
	}
	return false, nil
}
func (s *Store) BicycleAlreadyExists(id string) (bool, error) {
	var bicID string
	row, err := s.db.Query("select id from bicycles where id=?", id)
	if err != nil {
		return true, err
	}
	for row.Next() {

		err = row.Scan(&bicID)
		if err != nil {
			return true, err
		}
	}
	if bicID != "" || bicID == id {
		return true, nil
	}
	return false, nil
}

func (s *Store) UpdateAccessory(column string, value any, accessoryID string) error {
	log.Println("Entered function , column is : ", column, " , value is : ", value, "id of accessory : ", accessoryID)
	query := fmt.Sprintf("update accessories set %s = ? where id=? ", column)
	_, err := s.db.Exec(query, value, accessoryID)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) UpdateBicycle(column string, value any, accessoryID string) error {
	log.Println("Entered function , column is : ", column, " , value is : ", value, "id of bicycle : ", accessoryID)
	query := fmt.Sprintf("update bicycles set %s = ? where id=? ", column)
	_, err := s.db.Exec(query, value, accessoryID)
	if err != nil {
		return err
	}
	return nil
}
