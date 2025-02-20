package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/GakuruUri/lang-portal/internal/handlers"
	"github.com/GakuruUri/lang-portal/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	dbPath := filepath.Join(wd, "words.db")
	db, err := service.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Setup router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		// Dashboard endpoints
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/quick_stats", h.GetQuickStats)
			dashboard.GET("/study_progress", h.GetStudyProgress)
			dashboard.GET("/last_study_session", h.GetLastStudySession)
		}

		// Words endpoints
		api.GET("/words", h.GetWords)
		api.GET("/words/:id", h.GetWordByID)

		// Study sessions endpoints
		studySessions := api.Group("/study_sessions/:id")
		{
			studySessions.POST("/words/:word_id/review", h.CreateWordReview)
		}

		// Reset endpoints
		api.POST("/reset_history", h.ResetHistory)
		api.POST("/full_reset", h.FullReset)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
