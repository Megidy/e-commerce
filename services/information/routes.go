package information

import (
	"fmt"
	"log"

	templates "github.com/Megidy/e-commerce/frontend/templates/information"
	"github.com/Megidy/e-commerce/services/auth"
	"github.com/Megidy/e-commerce/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	templates        types.Templates
	userStore        types.UserStore
	informationStore types.InformationStore
}

func NewHandler(templates types.Templates, userStore types.UserStore, informationStore types.InformationStore) *Handler {
	return &Handler{
		templates:        templates,
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
	log.Println("user :", user)
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
	log.Println("question: ", question)
	log.Println("responds: ", responds)
	templates.LoadSingleQuestion(question, responds, user.Role).Render(c.Request.Context(), c.Writer)
}
func (h *Handler) LoadCreateQuestion(c *gin.Context) {
	templates.LoadCreateQuestionPage().Render(c.Request.Context(), c.Writer)
}
func (h *Handler) CreateQuestion(c *gin.Context) {
	var question types.Question

	u, ok := c.Get("user")
	if !ok {
		log.Println("user not found")
		return
	}
	user := u.(types.User)
	question.Id = uuid.NewString()
	question.UserID = user.ID
	question.Title = h.templates.GetDataFromForm(c, "title")
	question.Body = h.templates.GetDataFromForm(c, "body")
	err := h.informationStore.CreateQuestion(question)
	if err != nil {
		log.Println(err)
		return
	}
	c.Writer.Header().Add("HX-Redirect", fmt.Sprintf("/questions/%s", question.Id))
}

func (h *Handler) CreateRespond(c *gin.Context) {
	var respond types.Respond
	respond.QuestionID = c.Param("questionID")
	respond.RespondID = uuid.NewString()
	respond.Body = h.templates.GetDataFromForm(c, "respond")
	err := h.informationStore.CreateRespond(respond)
	if err != nil {
		log.Println(err)
		return
	}
	c.Writer.Header().Add("HX-Redirect", fmt.Sprintf("/questions/%s", respond.QuestionID))
}
func (h *Handler) DeleteQuestion(c *gin.Context) {

}
