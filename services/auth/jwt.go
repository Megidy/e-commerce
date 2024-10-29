package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Megidy/e-commerce/config"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	UserStore types.UserStore
}

func NewHandler(userStore types.UserStore) *Handler {
	return &Handler{
		UserStore: userStore,
	}
}

func (h *Handler) WithJWT(c *gin.Context) {

	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		AccesDenied(c)
	}
	token, err := ValidateJWT(tokenString)
	if err != nil {
		AccesDenied(c)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		id := claims["userID"].(string)
		user, err := h.UserStore.GetUserById(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

		}
		c.Set("user", user)
		c.Next()
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf(" Unexpected signing method: %v", token.Header["alg"])
		}
		config := config.InitConfig()
		return []byte(config.Secret), nil
	})
}

func CreateJWT(secret []byte, userId string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(time.Second * 60 * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AccesDenied(c *gin.Context) {
	c.AbortWithStatus(http.StatusUnauthorized)
}
