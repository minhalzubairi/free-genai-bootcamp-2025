package handlers

import (
    "log"
    "net/http"
    "strconv"
    "database/sql"

    "github.com/gin-gonic/gin"
    "github.com/minhalzubairi/lang-portal/backend-go/internal/service"
)

// GetGroups handles the GET /api/groups endpoint
func GetGroups(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    groups, pagination, err := service.GetGroups(page, perPage)
    if err != nil {
        log.Printf("Error getting groups: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "GROUPS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      groups,
        "pagination": pagination,
    })
}

// GetGroup handles the GET /api/groups/:id endpoint
func GetGroup(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid group ID",
            "code":  "INVALID_GROUP_ID",
        })
        return
    }

    group, err := service.GetGroup(id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Group not found",
                "code":  "GROUP_NOT_FOUND",
            })
            return
        }
        log.Printf("Error getting group %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "GROUP_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, group)
}

// GetGroupWords handles the GET /api/groups/:id/words endpoint
func GetGroupWords(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid group ID",
            "code":  "INVALID_GROUP_ID",
        })
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    words, pagination, err := service.GetGroupWords(id, page, perPage)
    if err != nil {
        log.Printf("Error getting words for group %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "GROUP_WORDS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      words,
        "pagination": pagination,
    })
}

// GetGroupStudySessions handles the GET /api/groups/:id/study_sessions endpoint
func GetGroupStudySessions(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid group ID",
            "code":  "INVALID_GROUP_ID",
        })
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    sessions, pagination, err := service.GetGroupStudySessions(id, page, perPage)
    if err != nil {
        log.Printf("Error getting study sessions for group %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "GROUP_SESSIONS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      sessions,
        "pagination": pagination,
    })
}
