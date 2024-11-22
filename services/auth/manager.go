package auth

import (
	"log"
	"net/http"

	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
)

type Manager struct {
	userStore types.UserStore
}

func NewManager(userStore types.UserStore) *Manager {
	return &Manager{userStore: userStore}
}

func (m *Manager) WithManagerRole(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		log.Println("no user found when trying to access role manager ")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	user := u.(types.User)
	if user.Role != "manager" {
		log.Println("access denied when trying to enter without role manager")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Next()
}
