package models

import (
	"time"
)

type Word struct {
	ID       int64  `json:"id"`
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
	English  string `json:"english"`
	Parts    string `json:"parts"` // JSON string
}

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type StudyActivity struct {
	ID              int64     `json:"id"`
	StudySessionID  int64     `json:"study_session_id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	StudySessionID  int64     `json:"study_session_id"`
	Correct         bool      `json:"correct"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
}

type WordReviewItem struct {
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

// Response structures matching the API spec
type DashboardLastStudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
	GroupName       string    `json:"group_name"`
}

type DashboardStudyProgress struct {
	TotalWordsStudied    int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

type DashboardQuickStats struct {
	SuccessRate       float64 `json:"success_rate"`
	TotalStudySessions int    `json:"total_study_sessions"`
	TotalActiveGroups  int    `json:"total_active_groups"`
	StudyStreakDays    int    `json:"study_streak_days"`
}
