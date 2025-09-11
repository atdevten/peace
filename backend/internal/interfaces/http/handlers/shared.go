package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard API Response format
type APIResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Response codes
const (
	CodeSuccess      = "SUCCESS"
	CodeBadRequest   = "BAD_REQUEST"
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeServerError  = "SERVER_ERROR"
)

// Success response with data
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error response with custom code and message
func Error(c *gin.Context, code string, message string) {
	status := http.StatusInternalServerError
	switch code {
	case CodeBadRequest:
		status = http.StatusBadRequest
	case CodeUnauthorized:
		status = http.StatusUnauthorized
	case CodeForbidden:
		status = http.StatusForbidden
	case CodeNotFound:
		status = http.StatusNotFound
	case CodeServerError:
		status = http.StatusInternalServerError
	default:
		status = http.StatusBadRequest
	}

	c.JSON(status, APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
