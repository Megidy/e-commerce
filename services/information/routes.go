package information

import (
	"log"

	templates "github.com/Megidy/e-commerce/frontend/templates/information"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userStore        types.UserStore
	informationStore types.InformationStore
}

func NewHandler(userStore types.UserStore, informationStore types.InformationStore) *Handler {
	return &Handler{
		userStore:        userStore,
		informationStore: informationStore,
	}
}

func (h *Handler) RegisterRoutes(router gin.IRouter, auth *auth.Handler, manager *auth.Manager) {
	router.GET("/information", h.InformationPage)
	router.GET("/questions", h.LoadQuestionsPage)
	router.GET("/questions/:questionID", auth.WithJWT, h.LoadSingleQuestion)
	router.GET("/questions/create", auth.WithJWT, h.LoadCreateQuestion)
	router.POST("/questions/create/confirm", auth.WithJWT, h.CreateQuestion)
	router.POST("/questions/:questionID/respond/create", auth.WithJWT, manager.WithManagerRole, h.CreateRespond)
	router.DELETE("/questions/:questionID/delete", auth.WithJWT, manager.WithManagerRole, h.DeleteQuestion)

}
func (h *Handler) InformationPage(c *gin.Context) {
	templates.LoadOverAllPage().Render(c.Request.Context(), c.Writer)
}

func (h *Handler) LoadQuestionsPage(c *gin.Context) {
	questions, err := h.informationStore.GetAllQuestions()
	if err != nil {
		log.Println(err)
		return
	}
	templates.LoadQuestionPage(questions).Render(c.Request.Context(), c.Writer)
}

func (h *Handler) LoadSingleQuestion(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)
	questionID := c.Param("questionID")

	question, err := h.informationStore.GetSingleQuestion(questionID)
	if err != nil {
		log.Println(err)
		return
	}
	responds, err := h.informationStore.GetAllResponds(questionID)
	if err != nil {
		log.Println(err)
		return
	}

	templates.LoadSingleQuestion(question, responds, user.Role)
}
func (h *Handler) LoadCreateQuestion(c *gin.Context) {

}
func (h *Handler) CreateQuestion(c *gin.Context) {

}

func (h *Handler) CreateRespond(c *gin.Context) {

}
func (h *Handler) DeleteQuestion(c *gin.Context) {

}
