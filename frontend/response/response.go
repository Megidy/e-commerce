package response

import (
	"github.com/gin-gonic/gin"
)

type ResponseHandler struct {
}

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

func (h *ResponseHandler) LoadResponse(c *gin.Context, statusCode int, file string, response map[string]any) {
	c.HTML(statusCode, file, response)

}
