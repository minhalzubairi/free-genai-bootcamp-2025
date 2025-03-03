package handlers

import (
    "database/sql"
    "log"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/minhalzubairi/lang-portal/backend-go/internal/service"
)

// GetLastStudySession handles the GET /api/dashboard/last_study_session endpoint
func GetLastStudySession(c *gin.Context) {
    session, err := service.GetLastStudySession()
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "No study sessions found",
                "code":  "NO_STUDY_SESSIONS",
            })
            return
        }
        log.Printf("Error getting last study session: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "LAST_SESSION_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, session)
}

// GetStudyProgress handles the GET /api/dashboard/study_progress endpoint
func GetStudyProgress(c *gin.Context) {
    progress, err := service.GetStudyProgress()
    if err != nil {
        log.Printf("Error getting study progress: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "STUDY_PROGRESS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, progress)
}

// GetQuickStats handles the GET /api/dashboard/quick-stats endpoint
func GetQuickStats(c *gin.Context) {
    stats, err := service.GetQuickStats()
    if err != nil {
        log.Printf("Error getting quick stats: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "QUICK_STATS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, stats)
}

// GetStudyActivities handles the GET /api/study_activities endpoint
func GetStudyActivities(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    activities, pagination, err := service.GetStudyActivities(page, perPage)
    if err != nil {
        log.Printf("Error getting study activities: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "STUDY_ACTIVITIES_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      activities,
        "pagination": pagination,
    })
}









// GetStudyActivity handles the GET /api/study_activities/:id endpoint
func GetStudyActivity(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid activity ID",
            "code":  "INVALID_ACTIVITY_ID",
        })
        return
    }

    activity, err := service.GetStudyActivity(id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Activity not found",
                "code":  "ACTIVITY_NOT_FOUND",
            })
            return
        }
        log.Printf("Error getting activity %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "ACTIVITY_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, activity)
}
// GetStudyActivitySessions handles the GET /api/study_activities/:id/study_sessions endpoint
func GetStudyActivitySessions(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid activity ID",
            "code":  "INVALID_ACTIVITY_ID",
        })
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    sessions, pagination, err := service.GetStudyActivitySessions(id, page, perPage)
    if err != nil {
        log.Printf("Error getting study sessions for activity %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "ACTIVITY_SESSIONS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      sessions,
        "pagination": pagination,
    })
}
// CreateStudyActivity handles the POST /api/study_activities endpoint
func CreateStudyActivity(c *gin.Context) {
    var req service.CreateActivityRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "code":  "INVALID_REQUEST",
            "details": err.Error(),
        })
        return
    }

    activity, err := service.CreateStudyActivity(&req)
    if err != nil {
        log.Printf("Error creating study activity: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "ACTIVITY_CREATE_ERROR",
        })
        return
    }

    c.JSON(http.StatusCreated, activity)
}
// GetStudySession handles the GET /api/study_sessions/:id endpoint
func GetStudySession(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid session ID",
            "code":  "INVALID_SESSION_ID",
        })
        return
    }

    session, err := service.GetStudySession(id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Study session not found",
                "code":  "SESSION_NOT_FOUND",
            })
            return
        }
        log.Printf("Error getting study session %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "SESSION_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, session)
}
// GetStudySessionWords handles the GET /api/study_sessions/:id/words endpoint
func GetStudySessionWords(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid session ID",
            "code":  "INVALID_SESSION_ID",
        })
        return
    }

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    words, pagination, err := service.GetStudySessionWords(id, page, perPage)
    if err != nil {
        log.Printf("Error getting words for session %d: %v", id, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "SESSION_WORDS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      words,
        "pagination": pagination,
    })
}
// GetStudySessions handles the GET /api/study_sessions endpoint
func GetStudySessions(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "100"))

    sessions, pagination, err := service.GetStudySessions(page, perPage)
    if err != nil {
        log.Printf("Error getting study sessions: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "STUDY_SESSIONS_FETCH_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items":      sessions,
        "pagination": pagination,
    })
}
// CreateWordReview handles the POST /api/study_sessions/:id/words/:word_id/review endpoint
func CreateWordReview(c *gin.Context) {
    sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid session ID",
            "code":  "INVALID_SESSION_ID",
        })
        return
    }

    wordID, err := strconv.ParseInt(c.Param("word_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid word ID",
            "code":  "INVALID_WORD_ID",
        })
        return
    }

    var req service.CreateWordReviewRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
            "code":  "INVALID_REQUEST",
            "details": err.Error(),
        })
        return
    }

    review, err := service.CreateWordReview(sessionID, wordID, &req)
    if err != nil {
        if err.Error() == "session not found" {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Study session not found",
                "code":  "SESSION_NOT_FOUND",
            })
            return
        }
        if err.Error() == "word not found" {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Word not found",
                "code":  "WORD_NOT_FOUND",
            })
            return
        }
        log.Printf("Error creating word review: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "WORD_REVIEW_CREATE_ERROR",
        })
        return
    }

    c.JSON(http.StatusCreated, review)
}
// ResetHistory handles the POST /api/reset_history endpoint
func ResetHistory(c *gin.Context) {
    if err := service.ResetHistory(); err != nil {
        log.Printf("Error resetting history: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "RESET_HISTORY_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Study history has been reset successfully",
    })
}
// FullReset handles the POST /api/full_reset endpoint
func FullReset(c *gin.Context) {
    if err := service.FullReset(); err != nil {
        log.Printf("Error performing full reset: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
            "code":  "FULL_RESET_ERROR",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Database has been fully reset successfully",
    })
}
