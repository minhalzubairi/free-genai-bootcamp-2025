package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/minhalzubairi/lang-portal/backend-go/internal/handlers"
	"github.com/minhalzubairi/lang-portal/backend-go/internal/service"
)

func main() {
	// Initialize database
	if err := service.InitDB("words.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer service.CloseDB()

	r := gin.Default()

	// API routes group
	api := r.Group("/api")
	{
		// Dashboard routes
		api.GET("/dashboard/last_study_session", handlers.GetLastStudySession)
		api.GET("/dashboard/study_progress", handlers.GetStudyProgress)
		api.GET("/dashboard/quick-stats", handlers.GetQuickStats)

		// Study activities routes
		api.GET("/study_activities", handlers.GetStudyActivities)
		api.GET("/study_activities/:id", handlers.GetStudyActivity)
		api.GET("/study_activities/:id/study_sessions", handlers.GetStudyActivitySessions)
		api.POST("/study_activities", handlers.CreateStudyActivity)

		// Words routes
		api.GET("/words", handlers.GetWords)
		api.GET("/words/:id", handlers.GetWord)

		// Groups routes
		api.GET("/groups", handlers.GetGroups)
		api.GET("/groups/:id", handlers.GetGroup)
		api.GET("/groups/:id/words", handlers.GetGroupWords)
		api.GET("/groups/:id/study_sessions", handlers.GetGroupStudySessions)

		// Study sessions routes
		api.GET("/study_sessions", handlers.GetStudySessions)
		api.GET("/study_sessions/:id", handlers.GetStudySession)
		api.GET("/study_sessions/:id/words", handlers.GetStudySessionWords)
		api.POST("/study_sessions/:id/words/:word_id/review", handlers.CreateWordReview)

		// Reset routes
		api.POST("/reset_history", handlers.ResetHistory)
		api.POST("/full_reset", handlers.FullReset)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 
