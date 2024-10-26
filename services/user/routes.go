package user

import (
	"log"
	"net/http"

	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/Megidy/e-commerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	userStore types.UserStore
}

func NewHandler(userStore types.UserStore) *Handler {
	return &Handler{
		userStore: userStore,
	}
}
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.Any("/signup", h.SignUp)

}
func (h *Handler) SignUp(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "signup.html", nil)
	} else if c.Request.Method == http.MethodPost {
		var payload types.CreateUserPayload
		var user types.User
		err := c.ShouldBindJSON(&payload)
		if err != nil {
			utils.HandleError(c, err, "failed to read", http.StatusBadRequest)
			return
		}
		log.Println(payload)
		hashedPassword, err := auth.HashPassword(payload.Password)
		if err != nil {
			utils.HandleError(c, err, "failed to hash password", http.StatusBadRequest)
			return
		}
		user = types.User{
			ID:       uuid.NewString(),
			Name:     payload.Name,
			LastName: payload.LastName,
			Email:    payload.Email,
			Password: string(hashedPassword),
		}
		ok, err := h.userStore.AlreadyExists(&user)
		if err != nil {
			utils.HandleError(c, err, "", http.StatusBadRequest)
			return

		}
		if ok {

			utils.HandleError(c, nil, "email is already been used", http.StatusBadRequest)
			return

		}
		err = h.userStore.CreateUser(&user)
		if err != nil {
			utils.HandleError(c, err, "", http.StatusInternalServerError)
			return
		}

	}
	c.Redirect(http.StatusFound, "/products")
}
