package service

import (
	"database/sql"

	"github.com/GakuruUri/lang-portal/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetQuickStats() (*models.QuickStats, error) {
	stats := &models.QuickStats{}

	// Get success rate
	err := db.QueryRow(`
		WITH review_stats AS (
			SELECT 
				COUNT(*) as total_reviews,
				SUM(CASE WHEN correct THEN 1 ELSE 0 END) as correct_reviews
			FROM word_review_items
		)
		SELECT 
			CASE 
				WHEN total_reviews > 0 
				THEN (CAST(correct_reviews AS FLOAT) / CAST(total_reviews AS FLOAT)) * 100 
				ELSE 0 
			END
		FROM review_stats
	`).Scan(&stats.SuccessRate)
	if err != nil {
		return nil, err
	}

	// Get total study sessions
	err = db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&stats.TotalStudySessions)
	if err != nil {
		return nil, err
	}

	// Get total active groups (groups with at least one study session)
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT group_id) 
		FROM study_sessions
	`).Scan(&stats.TotalActiveGroups)
	if err != nil {
		return nil, err
	}

	// Calculate study streak
	err = db.QueryRow(`
		WITH daily_sessions AS (
			SELECT DATE(created_at) as study_date
			FROM word_review_items
			GROUP BY DATE(created_at)
			ORDER BY study_date DESC
		),
		streak_calc AS (
			SELECT 
				study_date,
				ROW_NUMBER() OVER (ORDER BY study_date DESC) as row_num,
				JULIANDAY(study_date) as julian_date
			FROM daily_sessions
		)
		SELECT COUNT(*)
		FROM streak_calc
		WHERE julian_date = (
			SELECT MIN(julian_date) + row_num - 1
			FROM streak_calc
		)
	`).Scan(&stats.StudyStreakDays)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (db *DB) GetStudyProgress() (*models.StudyProgress, error) {
	progress := &models.StudyProgress{}

	err := db.QueryRow(`
		SELECT 
			(SELECT COUNT(DISTINCT word_id) FROM word_review_items) as studied,
			(SELECT COUNT(*) FROM words) as total
	`).Scan(&progress.TotalWordsStudied, &progress.TotalAvailableWords)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (db *DB) GetLastStudySession() (*models.StudySession, error) {
	session := &models.StudySession{}

	err := db.QueryRow(`
		SELECT 
			s.id,
			s.group_id,
			s.study_activity_id,
			s.start_time,
			g.name as group_name
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		ORDER BY s.start_time DESC
		LIMIT 1
	`).Scan(
		&session.ID,
		&session.GroupID,
		&session.StudyActivityID,
		&session.StartTime,
		&session.GroupName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return session, nil
}

func (db *DB) GetWords(page, itemsPerPage int) (*models.PaginatedResponse, error) {
	offset := (page - 1) * itemsPerPage

	// Get total count
	var totalItems int
	err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalItems)
	if err != nil {
		return nil, err
	}

	// Get words for current page
	rows, err := db.Query(`
		SELECT 
			w.id,
			w.japanese,
			w.romaji,
			w.english,
			w.parts,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(CASE WHEN NOT wri.correct THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id
		LIMIT ? OFFSET ?
	`, itemsPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var w models.WordWithStats
		err := rows.Scan(
			&w.ID,
			&w.Japanese,
			&w.Romaji,
			&w.English,
			&w.Parts,
			&w.CorrectCount,
			&w.WrongCount,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, w)
	}

	return &models.PaginatedResponse{
		Items: words,
		Pagination: models.Pagination{
			CurrentPage:  page,
			TotalPages:   (totalItems + itemsPerPage - 1) / itemsPerPage,
			TotalItems:   totalItems,
			ItemsPerPage: itemsPerPage,
		},
	}, nil
}

func (db *DB) GetWordByID(id int64) (*models.WordWithStats, error) {
	word := &models.WordWithStats{}

	err := db.QueryRow(`
		SELECT 
			w.id,
			w.japanese,
			w.romaji,
			w.english,
			w.parts,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(CASE WHEN NOT wri.correct THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`, id).Scan(
		&word.ID,
		&word.Japanese,
		&word.Romaji,
		&word.English,
		&word.Parts,
		&word.CorrectCount,
		&word.WrongCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get groups for this word
	rows, err := db.Query(`
		SELECT g.id, g.name
		FROM groups g
		JOIN word_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var g models.Group
		err := rows.Scan(&g.ID, &g.Name)
		if err != nil {
			return nil, err
		}
		word.Groups = append(word.Groups, g)
	}

	return word, nil
}

func (db *DB) CreateWordReview(sessionID, wordID int64, correct bool) error {
	_, err := db.Exec(`
		INSERT INTO word_review_items (word_id, study_session_id, correct)
		VALUES (?, ?, ?)
	`, wordID, sessionID, correct)
	return err
}

func (db *DB) ResetHistory() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (db *DB) FullReset() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	tables := []string{
		"word_review_items",
		"study_sessions",
		"word_groups",
		"words",
		"groups",
		"study_activities",
	}

	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM " + table)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
