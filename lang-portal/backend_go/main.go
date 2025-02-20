package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/GakuruUri/lang-portal/db"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

func main() {
	// Initialize database
	var err error
	dbConn, err = sql.Open("sqlite3", "./words.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Initialize router
	r := gin.Default()

	// Dashboard endpoints
	r.GET("/api/dashboard/last_study_session", getLastStudySession)
	r.GET("/api/dashboard/study_progress", getStudyProgress)
	r.GET("/api/dashboard/quick_stats", getQuickStats)

	// Study activities endpoints
	r.GET("/api/study_activities/:id", getStudyActivity)
	r.GET("/api/study-activities/:id/study_sessions", getStudyActivitySessions)
	r.POST("/api/study_activities", createStudyActivity)

	// Words endpoints
	r.GET("/api/words", getWords)
	r.GET("/api/words/:id", getWord)

	// Start server
	r.Run(":8080")
}

// Handler functions
func getLastStudySession(c *gin.Context) {
	session, err := db.GetLastStudySession(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func getStudyProgress(c *gin.Context) {
	progress, err := db.GetStudyProgress(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

func getQuickStats(c *gin.Context) {
	stats, err := db.GetQuickStats(dbConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func getStudyActivity(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}

func getStudyActivitySessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Not implemented yet",
	})
}

func createStudyActivity(c *gin.Context) {
	var request struct {
		GroupID int64 `json:"group_id" binding:"required"`
		StudyActivityID int64 `json:"study_activity_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activity, err := db.CreateStudyActivity(dbConn, request.GroupID, request.StudyActivityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

func getWords(c *gin.Context) {
	page := 1
	itemsPerPage := 100

	words, pagination, err := db.GetWords(dbConn, page, itemsPerPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": words,
		"pagination": pagination,
	})
}

func getWord(c *gin.Context) {
	id := c.Param("id")
	wordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	word, err := db.GetWord(dbConn, wordID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
}
