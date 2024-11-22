package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Megidy/e-commerce/config"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockUserStore struct {
}

func (m *mockUserStore) CreateUser(user *types.User) error {
	return nil
}
func (m *mockUserStore) AlreadyExists(user *types.User) (bool, error) {
	return true, nil
}
func (m *mockUserStore) GetUserByEmail(email string) (types.User, error) {

	return types.User{}, nil
}

func (m *mockUserStore) GetUserById(id string) (types.User, error) {
	if id == "123 plz pass" {
		return types.User{
			ID:       "123 plz pass",
			Name:     "testNameuserauthoorization",
			LastName: "testLastName",
			Email:    "email@123pleasypass",
			Password: "testpass",
			Created:  "11.11.1111",
			Role:     "user",
		}, nil
	}
	return types.User{}, nil
}

func TestCreateJWT(t *testing.T) {
	secret := "superpupersecretkey"
	userid := uuid.NewString()
	token, err := CreateJWT([]byte(secret), userid)

	if err != nil {
		t.Errorf("error while creating token : %v", err)
	}
	if token == "" {
		t.Errorf("expected token to be not empty , error: %v", token)
	}

}
func TestJWT(t *testing.T) {
	userStore := &mockUserStore{}
	t.Run("Should Pass if user is authorized", func(t *testing.T) {
		JWTauth := NewJWT(userStore)
		var u types.User = types.User{
			ID:       "123 plz pass",
			Name:     "testNameuserauthoorization",
			LastName: "testLastName",
			Email:    "email@123pleasypass",
			Password: "testpass",
			Created:  "11.11.1111",
			Role:     "user",
		}
		req, err := http.NewRequest(http.MethodGet, "/user", nil)
		if err != nil {
			t.Errorf("error while creating new request : %v", req)
		}
		config := config.InitConfig()

		secret, _ := CreateJWT([]byte(config.Secret), u.ID)
		req.AddCookie(&http.Cookie{
			Name:  "Authorization",
			Value: secret,
		})
		rr := httptest.NewRecorder()
		router := gin.Default()

		router.GET("/user", func(c *gin.Context) {

			secret, _ := CreateJWT([]byte(config.Secret), u.ID)
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("Authorization", secret, 3600*24*10, "/", "", false, true)

			JWTauth.WithJWT(c)
			if h := c.Writer.Header().Get("authorization"); h == "false" {
				t.Errorf("expected header to be false : %s", h)
			}
			cookie, _ := c.Cookie("Authorization")
			if cookie == "" {
				t.Errorf("expected cookie to be not null , cookie : %s", cookie)
			}

		})
		router.ServeHTTP(rr, req)
	})
}
