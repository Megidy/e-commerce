package user

import (
	"log"
	"net/http"

	"github.com/Megidy/e-commerce/config"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	response  types.LoadResponse
	userStore types.UserStore
}

func NewHandler(response types.LoadResponse, userStore types.UserStore) *Handler {
	return &Handler{
		response:  response,
		userStore: userStore,
	}
}
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.Any("/signup", h.SignUp)
	router.Any("/login", h.LogIn)
}
func (h *Handler) SignUp(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		h.response.LoadResponse(c, http.StatusOK, "signup.html", nil)
		return
	}

	if c.Request.Method == http.MethodPost {
		var payload types.CreateUserPayload
		var user types.User
		err := c.ShouldBindJSON(&payload)
		if err != nil {
			h.response.LoadResponse(c, http.StatusBadRequest, "signup.html", gin.H{
				"Error": err.Error(),
			})
			return
			// utils.HandleError(c, err, "failed to read", http.StatusBadRequest)
			// return
		}
		log.Println(payload)
		hashedPassword, err := auth.HashPassword(payload.Password)
		if err != nil {
			h.response.LoadResponse(c, http.StatusBadRequest, "signup.html", gin.H{
				"Error": err.Error(),
			})
			return
			// utils.HandleError(c, err, "failed to hash password", http.StatusBadRequest)
			// return
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
			h.response.LoadResponse(c, http.StatusBadRequest, "signup.html", gin.H{
				"Error": err.Error(),
			})
			return
			// utils.HandleError(c, err, "", http.StatusBadRequest)
			// return
		}

		if ok {
			h.response.LoadResponse(c, http.StatusBadRequest, "signup.html", gin.H{
				"Message": "Email is already in use",
			})
			return
		}

		err = h.userStore.CreateUser(&user)
		if err != nil {
			h.response.LoadResponse(c, http.StatusBadRequest, "signup.html", gin.H{
				"Error": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusSeeOther, "/login")
	}
}

func (h *Handler) LogIn(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		h.response.LoadResponse(c, http.StatusOK, "login.html", nil)
	} else if c.Request.Method == http.MethodPost {
		var logInPayload types.LogInPayload
		c.ShouldBindJSON(&logInPayload)
		ok, err := h.userStore.AlreadyExists(&types.User{Email: logInPayload.Email})
		if err != nil {
			h.response.LoadResponse(c, http.StatusBadRequest, "login.html", gin.H{
				"Error": err.Error(),
			})
			return
		}
		if !ok {
			h.response.LoadResponse(c, http.StatusBadRequest, "login.html", gin.H{
				"Message": "no user with this email found ",
			})
			return
		} else if ok {
			user, err := h.userStore.GetUserByEmail(logInPayload.Email)
			if err != nil {
				h.response.LoadResponse(c, http.StatusBadRequest, "login.html", gin.H{
					"Error": err.Error(),
				})
				return
			}
			ok := auth.ComparePassword(user.Password, logInPayload.Password)
			if !ok {
				h.response.LoadResponse(c, http.StatusBadRequest, "login.html", gin.H{
					"Error": "Invlid data sent",
				})
				return
			}

			config := config.InitConfig()
			secret, err := auth.CreateJWT([]byte(config.Secret), user.ID)
			if err != nil {
				h.response.LoadResponse(c, http.StatusBadRequest, "login.html", gin.H{
					"Error": err.Error(),
				})
				return
			}
			log.Println(secret)
			log.Println("user logged in")
		}
	}

}
