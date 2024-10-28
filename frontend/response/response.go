package response

import (
	"log"
	"text/template"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{}
}

func (h *TemplateHandler) LoadResponse(c *gin.Context, statusCode int, file string, response map[string]any) {
	log.Println(response)
	c.HTML(statusCode, file, response)

}

func (h *TemplateHandler) ExecuteTemplate(c *gin.Context, filePath string, data map[string]any) error {
	tmpl := template.Must(template.ParseFiles(filePath))

	return tmpl.Execute(c.Writer, data)
}

func (h *TemplateHandler) ExecuteSpecificTemplate(c *gin.Context, name, filePath string, data map[string]any) error {
	tmpl := template.Must(template.ParseFiles(filePath))
	return tmpl.ExecuteTemplate(c.Writer, name, data)
}

func (h *TemplateHandler) GetDataFromForm(c *gin.Context, key string) string {
	return c.Request.PostFormValue(key)
}
