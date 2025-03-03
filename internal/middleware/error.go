package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorHandler middleware for consistent error handling
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: err.Error(),
			})
			return
		}
	}
}
