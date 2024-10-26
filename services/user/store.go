package user

import (
	"database/sql"

	"github.com/Megidy/e-commerce/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateUser(user *types.User) error {
	_, err := s.db.Exec("insert into users(id,name,lastname,email,password,role) values(?,?,?,?,?,?)",
		user.ID, user.Name, user.LastName, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) AlreadyExists(user *types.User) (bool, error) {
	var email string

	row := s.db.QueryRow("select email from users where email=?", user.Email)

	err := row.Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, err
	}
	if email == user.Email {
		return true, nil
	}
	return false, nil
}
func (s *Store) GetUserByEmail(email string) (types.User, error) {
	var u types.User
	row, err := s.db.Query("select * from users where email=?", email)
	if err != nil {
		return types.User{}, err
	}
	for row.Next() {
		err = row.Scan(&u.ID, &u.Name, &u.LastName, &u.Email, &u.Password, &u.Created, &u.Role)
		if err != nil {
			return types.User{}, err
		}
	}
	return u, nil

}
func (s *Store) GetUserById(ID string) (types.User, error) {
	var u types.User
	row, err := s.db.Query("select * from users where id=?", ID)
	if err != nil {
		return types.User{}, err
	}
	for row.Next() {
		err = row.Scan(&u.ID, &u.Name, &u.LastName, &u.Email, &u.Password, &u.Created, &u.Role)
		if err != nil {
			return types.User{}, err
		}
	}
	return u, nil

}
