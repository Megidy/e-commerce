package types

import "github.com/gin-gonic/gin"

type User struct {
	ID       string
	Name     string
	LastName string
	Email    string
	Password string
	Created  string
	Role     string
}

type CreateUserPayload struct {
	Name     string
	LastName string
	Email    string
	Password string
}

type LogInPayload struct {
	Email    string
	Password string
}

type UserStore interface {
	CreateUser(user *User) error
	AlreadyExists(user *User) (bool, error)
	GetUserByEmail(email string) (User, error)
	GetUserById(id string) (User, error)
}

type LoadResponse interface {
	LoadResponse(c *gin.Context, statusCode int, file string, response map[string]any)
}
