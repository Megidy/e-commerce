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
	_, err := s.db.Exec("insert into users(id,name,lastname,email,password) values(?,?,?,?,?)",
		user.ID, user.Name, user.LastName, user.Email, user.Password)
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
	return true, nil
}
