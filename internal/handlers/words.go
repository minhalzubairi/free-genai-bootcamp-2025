package handlers

import (
    "database/sql"
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

// GetWord handles the GET /api/words/:id endpoint
func GetWord(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid word ID",
            "code":  "INVALID_WORD_ID",
        })
        return
    }

    word, err := service.GetWord(id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Word not found",
                "code":  "WORD_NOT_FOUND",
            })
            return
        }
        log.Printf("Error getting word %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "WORD_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, word)
}
