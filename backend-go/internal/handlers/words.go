package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minhalzubairi/lang-portal/backend-go/internal/service"
)

// GetWords handles the GET /api/words endpoint
func GetWords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

	words, pagination, err := service.GetWords(page, perPage)
	if err != nil {
		log.Printf("Error getting words: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  "WORDS_FETCH_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":      words,
		"pagination": pagination,
	})
}
