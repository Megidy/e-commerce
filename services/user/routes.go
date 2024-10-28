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
	templates types.Templates
	userStore types.UserStore
}

func NewHandler(templates types.Templates, userStore types.UserStore) *Handler {
	return &Handler{
		templates: templates,
		userStore: userStore,
	}
}
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/signup", h.LoadSignUpHTML)
	router.POST("/signup/create", h.SignUp)
	router.GET("/login", h.LoadLogInHTML)
	router.POST("/login/enter", h.LogIn)

}

func (h *Handler) SignUp(c *gin.Context) {
	var payload types.CreateUserPayload
	var user types.User
	payload.Name = h.templates.GetDataFromForm(c, "name")
	payload.LastName = h.templates.GetDataFromForm(c, "lastname")
	payload.Email = h.templates.GetDataFromForm(c, "email")
	payload.Password = h.templates.GetDataFromForm(c, "password")
	log.Println(payload)
	ok, err := h.userStore.AlreadyExists(&types.User{Email: payload.Email})
	if err != nil {
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/signup.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}
	if ok {
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/signup.html", map[string]any{
			"Message": "Email is already taken",
		})
		return
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/signup.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}
	user = types.User{
		ID:       uuid.NewString(),
		Name:     payload.Name,
		LastName: payload.LastName,
		Email:    payload.Email,
		Password: string(hashedPassword),
		Role:     "User",
	}
	err = h.userStore.CreateUser(&user)
	if err != nil {
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/signup.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}

	h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/signup.html", map[string]any{
		"Message": "User Signed In",
	})
}

func (h *Handler) LoadSignUpHTML(c *gin.Context) {
	h.templates.ExecuteTemplate(c, "frontend/templates/signup.html", nil)
}
func (h *Handler) LoadLogInHTML(c *gin.Context) {
	h.templates.ExecuteTemplate(c, "frontend/templates/login.html", nil)
}

func (h *Handler) LogIn(c *gin.Context) {
	var logInPayload types.LogInPayload

	logInPayload.Email = h.templates.GetDataFromForm(c, "email")
	logInPayload.Password = h.templates.GetDataFromForm(c, "password")

	log.Println(logInPayload)

	ok, err := h.userStore.AlreadyExists(&types.User{Email: logInPayload.Email})
	if err != nil {
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/login.html", map[string]any{
			"Message": err.Error(),
		})
		return
	}
	if !ok {
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/login.html", map[string]any{
			"Message": "User not found",
		})
		return
	} else if ok {
		user, err := h.userStore.GetUserByEmail(logInPayload.Email)
		if err != nil {
			h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/login.html", map[string]any{
				"Message": "User not found",
			})
			return
		}

		ok := auth.ComparePassword(user.Password, logInPayload.Password)
		if !ok {
			h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/login.html", map[string]any{
				"Message": "Invalid data",
			})
			return
		}
		config := config.InitConfig()
		secret, err := auth.CreateJWT([]byte(config.Secret), user.ID)
		if err != nil {
			h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/login.html", map[string]any{
				"Message": err.Error(),
			})
			return
		}
		log.Println(secret)
		h.templates.ExecuteSpecificTemplate(c, "message-after-submit", "frontend/templates/login.html", map[string]any{
			"Message": "Successfuly logged in",
		})
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", secret, 3600*24*10, "", "", false, true)
	}
}
