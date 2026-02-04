package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Data    interface{} `json:"data"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
}

type PaginatedData struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Data:    data,
		Code:    CodeSuccess,
		Message: message,
	})
}

func SuccessList(c *gin.Context, items interface{}, total int64, page int, size int, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Data: PaginatedData{
			Items: items,
			Total: total,
			Page:  page,
			Size:  size,
		},
		Code:    CodeSuccess,
		Message: message,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Data:    data,
		Code:    CodeCreated,
		Message: message,
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, APIResponse{
		Data:    nil,
		Code:    code,
		Message: message,
	})
}
