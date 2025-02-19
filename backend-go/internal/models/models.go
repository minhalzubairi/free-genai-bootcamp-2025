package models

import (
	"database/sql"
	"time"
)

type Word struct {
	ID      int64          `json:"id"`
	Arabic  string         `json:"arabic"`
	Roman   string         `json:"roman"`
	English string         `json:"english"`
	Parts   sql.NullString `json:"parts,omitempty"`
}

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type StudySession struct {
	ID              int64     `json:"id"`
	GroupID         int64     `json:"group_id"`
	CreatedAt       time.Time `json:"created_at"`
	StudyActivityID int64     `json:"study_activity_id"`
}

type StudyActivity struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
	LaunchURL    string `json:"launch_url"`
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
