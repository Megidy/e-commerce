package response

import (
	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) GetDataFromForm(c *gin.Context, key string) string {
	return c.Request.PostFormValue(key)
}
