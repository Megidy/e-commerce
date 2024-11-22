package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestWithManagerRole(t *testing.T) {
	userStore := &mockUserStore{}
	m := NewManager(userStore)
	t.Run("should pass if users role is manager ", func(t *testing.T) {
		var mockedUser types.User = types.User{
			ID:       uuid.NewString(),
			Name:     "testName",
			LastName: "testLastName",
			Email:    "email@123pleasypass",
			Password: "testpass",
			Created:  "11.11.1111",
			Role:     "manager",
		}
		req, err := http.NewRequest(http.MethodGet, "/user", nil)
		if err != nil {
			t.Errorf("error while creating new request no /user : %v", err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/user", func(c *gin.Context) {
			c.Set("user", mockedUser)
			m.WithManagerRole(c)
		})
		router.ServeHTTP(rr, req)
		if rr.Code == http.StatusNotFound {
			t.Errorf("expected : %v, got :%v ", http.StatusOK, rr.Code)
		}
	})
	t.Run("Should Fail if user role is not manager ", func(t *testing.T) {
		var mockedUser types.User = types.User{
			ID:       uuid.NewString(),
			Name:     "testName",
			LastName: "testLastName",
			Email:    "email@123pleasypass",
			Password: "testpass",
			Created:  "11.11.1111",
			Role:     "user",
		}
		req, err := http.NewRequest(http.MethodGet, "/user", nil)
		if err != nil {
			t.Errorf("error while creating new request no /user : %v", err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.GET("/user", func(c *gin.Context) {
			c.Set("user", mockedUser)
			m.WithManagerRole(c)
		})
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected : %v, got :%v ", http.StatusNotFound, rr.Code)
		}
	})

}
