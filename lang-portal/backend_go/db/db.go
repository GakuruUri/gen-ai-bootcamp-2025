package db

import (
	"database/sql"
	"github.com/GakuruUri/lang-portal/models"
)

// GetLastStudySession returns the most recent study session
func GetLastStudySession(db *sql.DB) (*models.DashboardLastStudySession, error) {
	query := `
		SELECT 
			s.id, 
			s.group_id, 
			s.created_at, 
			s.study_activity_id,
			g.name as group_name
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT 1
	`
	
	var session models.DashboardLastStudySession
	err := db.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetStudyProgress returns the study progress statistics
func GetStudyProgress(db *sql.DB) (*models.DashboardStudyProgress, error) {
	query := `
		WITH studied_words AS (
			SELECT DISTINCT word_id 
			FROM word_reviews_items
		)
		SELECT 
			(SELECT COUNT(*) FROM studied_words) as total_words_studied,
			(SELECT COUNT(*) FROM words) as total_available_words
	`
	
	var progress models.DashboardStudyProgress
	err := db.QueryRow(query).Scan(
		&progress.TotalWordsStudied,
		&progress.TotalAvailableWords,
	)
	if err != nil {
		return nil, err
	}
	return &progress, nil
}

// GetQuickStats returns quick overview statistics
func GetQuickStats(db *sql.DB) (*models.DashboardQuickStats, error) {
	query := `
		WITH review_stats AS (
			SELECT 
				COUNT(*) as total_reviews,
				SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_reviews
			FROM word_reviews_items wri
			WHERE wri.created_at >= datetime('now', '-30 days')
		),
		session_stats AS (
			SELECT
				COUNT(DISTINCT ss.id) as total_sessions,
				COUNT(DISTINCT ss.group_id) as active_groups
			FROM study_sessions ss
			WHERE ss.created_at >= datetime('now', '-30 days')
		)
		SELECT 
			CASE 
				WHEN rs.total_reviews > 0 THEN 
					(CAST(rs.correct_reviews AS FLOAT) / rs.total_reviews) * 100
				ELSE 0
			END as success_rate,
			COALESCE(ss.total_sessions, 0) as total_sessions,
			COALESCE(ss.active_groups, 0) as active_groups,
			(
				SELECT COUNT(DISTINCT DATE(created_at))
				FROM study_sessions
				WHERE created_at >= datetime('now', '-30 days')
			) as streak_days
		FROM review_stats rs, session_stats ss
	`
	
	var stats models.DashboardQuickStats
	err := db.QueryRow(query).Scan(
		&stats.SuccessRate,
		&stats.TotalStudySessions,
		&stats.TotalActiveGroups,
		&stats.StudyStreakDays,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetWords returns a paginated list of words
func GetWords(db *sql.DB, page int, itemsPerPage int) ([]models.Word, *models.Pagination, error) {
	offset := (page - 1) * itemsPerPage
	
	// Get total count
	var totalItems int
	err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalItems)
	if err != nil {
		return nil, nil, err
	}
	
	// Get words for current page
	query := `
		SELECT id, japanese, romaji, english, parts
		FROM words
		LIMIT ? OFFSET ?
	`
	
	rows, err := db.Query(query, itemsPerPage, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	
	var words []models.Word
	for rows.Next() {
		var word models.Word
		err := rows.Scan(&word.ID, &word.Japanese, &word.Romaji, &word.English, &word.Parts)
		if err != nil {
			return nil, nil, err
		}
		words = append(words, word)
	}
	
	pagination := &models.Pagination{
		CurrentPage: page,
		ItemsPerPage: itemsPerPage,
		TotalItems: totalItems,
		TotalPages: (totalItems + itemsPerPage - 1) / itemsPerPage,
	}
	
	return words, pagination, nil
}

// GetWord returns a single word by ID
func GetWord(db *sql.DB, id int64) (*models.Word, error) {
	query := `
		SELECT id, japanese, romaji, english, parts
		FROM words
		WHERE id = ?
	`
	
	var word models.Word
	err := db.QueryRow(query, id).Scan(
		&word.ID,
		&word.Japanese,
		&word.Romaji,
		&word.English,
		&word.Parts,
	)
	if err != nil {
		return nil, err
	}
	return &word, nil
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(db *sql.DB, groupID int64, studyActivityID int64) (*models.StudyActivity, error) {
	query := `
		INSERT INTO study_activities (group_id, study_session_id)
		VALUES (?, ?)
		RETURNING id, group_id, study_session_id, created_at
	`
	
	var activity models.StudyActivity
	err := db.QueryRow(query, groupID, studyActivityID).Scan(
		&activity.ID,
		&activity.GroupID,
		&activity.StudySessionID,
		&activity.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &activity, nil
}
