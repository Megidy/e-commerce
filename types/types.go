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

type Templates interface {
	LoadResponse(c *gin.Context, statusCode int, file string, response map[string]any)
	ExecuteTemplate(c *gin.Context, filePath string, data map[string]any) error
	ExecuteSpecificTemplate(c *gin.Context, name, filePath string, data map[string]any) error
	GetDataFromForm(c *gin.Context, key string) string
}
