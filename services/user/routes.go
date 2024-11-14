package user

import (
	"log"
	"net/http"

	"github.com/Megidy/e-commerce/config"
	templates "github.com/Megidy/e-commerce/frontend/templates/user"
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
	router.GET("/", h.LoadSignUpTemplate)
	router.GET("/signup", h.LoadSignUpTemplate)
	router.POST("/signup/create", h.SignUp)
	router.GET("/login", h.LoadLogInTemplate)
	router.POST("/login/enter", h.LogIn)
	router.GET("/user", auth.NewHandler(h.userStore).WithJWT, h.UserAccount)

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
		templates.Signup(true, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	if ok {
		templates.Signup(true, "user already exists").Render(c.Request.Context(), c.Writer)
		return
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		templates.Signup(true, "failed to hash password").Render(c.Request.Context(), c.Writer)
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
		templates.Signup(true, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	c.Writer.Header().Add("HX-Redirect", "/login")
}

func (h *Handler) LoadSignUpTemplate(c *gin.Context) {
	templates.Signup(false, "").Render(c.Request.Context(), c.Writer)
}
func (h *Handler) LoadLogInTemplate(c *gin.Context) {
	templates.Login(false, "").Render(c.Request.Context(), c.Writer)
}

func (h *Handler) LogIn(c *gin.Context) {
	var logInPayload types.LogInPayload

	logInPayload.Email = h.templates.GetDataFromForm(c, "email")
	logInPayload.Password = h.templates.GetDataFromForm(c, "password")

	log.Println("login payload :", logInPayload)

	ok, err := h.userStore.AlreadyExists(&types.User{Email: logInPayload.Email})
	if err != nil {
		templates.Login(true, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	if !ok {
		c.Writer.Header().Add("hasEmail", "false")
		templates.Login(true, "Invalid data sent").Render(c.Request.Context(), c.Writer)
		return
	} else if ok {
		user, err := h.userStore.GetUserByEmail(logInPayload.Email)
		if err != nil {
			templates.Login(true, "Invalid data sent").Render(c.Request.Context(), c.Writer)
			return
		}

		ok := auth.ComparePassword(user.Password, logInPayload.Password)
		if !ok {
			c.Writer.Header().Add("correctPassword", "false")
			templates.Login(true, "Invalid data sent").Render(c.Request.Context(), c.Writer)
			return
		}
		config := config.InitConfig()
		secret, err := auth.CreateJWT([]byte(config.Secret), user.ID)
		if err != nil {
			templates.Login(true, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
		c.Writer.Header().Add("HX-Redirect", "/products/accessories")

		log.Println("cookie :", secret)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", secret, 3600*24*10, "", "", false, true)

	}
}
func (h *Handler) UserAccount(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)
	log.Println(user)
	templates.UserAccount(user).Render(c.Request.Context(), c.Writer)
}
