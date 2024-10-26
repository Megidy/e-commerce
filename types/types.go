package types

type User struct {
	ID       string
	Name     string
	LastName string
	Email    string
	Password string
	Created  string
}

type CreateUserPayload struct {
	Name     string
	LastName string
	Email    string
	Password string
}
type UserStore interface {
	CreateUser(user *User) error
	AlreadyExists(user *User) (bool, error)
}
