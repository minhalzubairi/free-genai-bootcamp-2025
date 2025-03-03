package service

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "github.com/minhalzubairi/lang-portal/backend-go/internal/models"
)

// StudySessionDetailResponse represents a detailed study session
type StudySessionDetailResponse struct {
    ID           int64            `json:"id"`
    ActivityName string           `json:"activity_name"`
    GroupName    string           `json:"group_name"`
    StartTime    time.Time        `json:"start_time"`
    Stats        StudySessionStats `json:"stats"`
    Words        []struct {
        Arabic     string    `json:"arabic"`
        Roman      string    `json:"roman"`
        English    string    `json:"english"`
        IsCorrect  bool      `json:"is_correct"`
        ReviewedAt time.Time `json:"reviewed_at"`
    } `json:"words"`
}

// StudySessionStats represents statistics for a study session
type StudySessionStats struct {
    TotalWords   int `json:"total_words"`
    CorrectCount int `json:"correct_count"`
    WrongCount   int `json:"wrong_count"`
}

// StudySessionResponse represents a study session with its details
type StudySessionResponse struct {
    ID           int64            `json:"id"`
    ActivityName string           `json:"activity_name"`
    GroupName    string           `json:"group_name"`
    StartTime    time.Time        `json:"start_time"`
    Stats        StudySessionStats `json:"stats"`
}

// GetStudySessions returns a paginated list of study sessions
func GetStudySessions(page, perPage int) ([]StudySessionResponse, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&total)
    if err != nil {
        log.Printf("Error counting study sessions: %v", err)
        return nil, nil, err
    }

    // Get sessions with their stats
    rows, err := db.Query(`
        SELECT 
            ss.id,
            sa.name as activity_name,
            g.name as group_name,
            ss.created_at as start_time,
            COUNT(DISTINCT wri.word_id) as total_words,
            SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END) as wrong_count
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        JOIN groups g ON ss.group_id = g.id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        GROUP BY ss.id
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?`,
        perPage, offset)
    if err != nil {
        log.Printf("Error querying study sessions: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var sessions []StudySessionResponse
    for rows.Next() {
        var s StudySessionResponse
        err := rows.Scan(
            &s.ID,
            &s.ActivityName,
            &s.GroupName,
            &s.StartTime,
            &s.Stats.TotalWords,
            &s.Stats.CorrectCount,
            &s.Stats.WrongCount,
        )
        if err != nil {
            log.Printf("Error scanning session: %v", err)
            return nil, nil, err
        }
        sessions = append(sessions, s)
    }

    pagination := &models.Pagination{
        CurrentPage:  page,
        ItemsPerPage: perPage,
        TotalItems:   total,
        TotalPages:   (total + perPage - 1) / perPage,
    }

    return sessions, pagination, nil
}

// GetStudySession returns a single study session by ID with its details
func GetStudySession(id int64) (*StudySessionDetailResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    var session StudySessionDetailResponse

    // Get session details
    err := db.QueryRow(`
        SELECT 
            ss.id,
            sa.name as activity_name,
            g.name as group_name,
            ss.created_at as start_time,
            COUNT(DISTINCT wri.word_id) as total_words,
            SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END) as wrong_count
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        JOIN groups g ON ss.group_id = g.id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        WHERE ss.id = ?
        GROUP BY ss.id`,
        id).Scan(
            &session.ID,
            &session.ActivityName,
            &session.GroupName,
            &session.StartTime,
            &session.Stats.TotalWords,
            &session.Stats.CorrectCount,
            &session.Stats.WrongCount,
        )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        log.Printf("Error getting study session %d: %v", id, err)
        return nil, err
    }

    // Get reviewed words
    rows, err := db.Query(`
        SELECT 
            w.arabic,
            w.roman,
            w.english,
            wri.correct as is_correct,
            wri.created_at as reviewed_at
        FROM word_review_items wri
        JOIN words w ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
        ORDER BY wri.created_at`,
        id)
    if err != nil {
        log.Printf("Error getting session words: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var word struct {
            Arabic     string    `json:"arabic"`
            Roman      string    `json:"roman"`
            English    string    `json:"english"`
            IsCorrect  bool      `json:"is_correct"`
            ReviewedAt time.Time `json:"reviewed_at"`
        }
        err := rows.Scan(
            &word.Arabic,
            &word.Roman,
            &word.English,
            &word.IsCorrect,
            &word.ReviewedAt,
        )
        if err != nil {
            log.Printf("Error scanning word: %v", err)
            return nil, err
        }
        session.Words = append(session.Words, word)
    }

    return &session, nil
}

// ResetHistory deletes all study sessions and word reviews
func ResetHistory() error {
    db := GetDB()
    if db == nil {
        return fmt.Errorf("database connection not initialized")
    }

    tx, err := db.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return err
    }

    // Delete all word reviews
    _, err = tx.Exec("DELETE FROM word_review_items")
    if err != nil {
        tx.Rollback()
        log.Printf("Error deleting word reviews: %v", err)
        return err
    }

    // Delete all study sessions
    _, err = tx.Exec("DELETE FROM study_sessions")
    if err != nil {
        tx.Rollback()
        log.Printf("Error deleting study sessions: %v", err)
        return err
    }

    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return err
    }

    return nil
}

// FullReset deletes all data and resets the database to initial state
func FullReset() error {
    db := GetDB()
    if db == nil {
        return fmt.Errorf("database connection not initialized")
    }

    tx, err := db.Begin()
    if err != nil {
        log.Printf("Error starting transaction: %v", err)
        return err
    }

    // Delete all data in reverse order of dependencies
    tables := []string{
        "word_review_items",
        "study_sessions",
        "words_groups",
        "words",
        "groups",
        "study_activities",
    }

    for _, table := range tables {
        _, err = tx.Exec("DELETE FROM " + table)
        if err != nil {
            tx.Rollback()
            log.Printf("Error deleting from %s: %v", table, err)
            return err
        }

        // Try to reset the autoincrement counter if possible
        _, err = tx.Exec("UPDATE SQLITE_SEQUENCE SET SEQ=0 WHERE NAME=?", table)
        if err != nil {
            // Ignore errors here as the table might not exist
            log.Printf("Note: Could not reset sequence for %s (this is usually OK)", table)
        }
    }

    if err := tx.Commit(); err != nil {
        log.Printf("Error committing transaction: %v", err)
        return err
    }

    return nil
}

// SessionWordResponse represents a word in a study session
type SessionWordResponse struct {
    Arabic     string    `json:"arabic"`
    Roman      string    `json:"roman"`
    English    string    `json:"english"`
    IsCorrect  bool      `json:"is_correct"`
    ReviewedAt time.Time `json:"reviewed_at"`
}

// CreateWordReviewRequest represents the request to create a word review
type CreateWordReviewRequest struct {
    IsCorrect bool `json:"is_correct" binding:"required"`
}

// CreateWordReviewResponse represents the response after creating a word review
type CreateWordReviewResponse struct {
    ID         int64     `json:"id"`
    WordID     int64     `json:"word_id"`
    SessionID  int64     `json:"session_id"`
    IsCorrect  bool      `json:"is_correct"`
    ReviewedAt time.Time `json:"reviewed_at"`
}

// GetStudySessionWords returns all words reviewed in a study session
func GetStudySessionWords(sessionID int64, page, perPage int) ([]SessionWordResponse, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow(`
        SELECT COUNT(*)
        FROM word_review_items
        WHERE study_session_id = ?`,
        sessionID).Scan(&total)
    if err != nil {
        log.Printf("Error counting session words: %v", err)
        return nil, nil, err
    }

    // Get words with their review status
    rows, err := db.Query(`
        SELECT 
            w.arabic,
            w.roman,
            w.english,
            wri.correct as is_correct,
            wri.created_at as reviewed_at
        FROM word_review_items wri
        JOIN words w ON wri.word_id = w.id
        WHERE wri.study_session_id = ?
        ORDER BY wri.created_at
        LIMIT ? OFFSET ?`,
        sessionID, perPage, offset)
    if err != nil {
        log.Printf("Error querying session words: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var words []SessionWordResponse
    for rows.Next() {
        var w SessionWordResponse
        err := rows.Scan(
            &w.Arabic,
            &w.Roman,
            &w.English,
            &w.IsCorrect,
            &w.ReviewedAt,
        )
        if err != nil {
            log.Printf("Error scanning word: %v", err)
            return nil, nil, err
        }
        words = append(words, w)
    }

    pagination := &models.Pagination{
        CurrentPage:  page,
        ItemsPerPage: perPage,
        TotalItems:   total,
        TotalPages:   (total + perPage - 1) / perPage,
    }

    return words, pagination, nil
}

// CreateWordReview creates a new word review for a study session
func CreateWordReview(sessionID, wordID int64, req *CreateWordReviewRequest) (*CreateWordReviewResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    // Verify session exists
    var sessionExists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM study_sessions WHERE id = ?)", sessionID).Scan(&sessionExists)
    if err != nil {
        log.Printf("Error checking session existence: %v", err)
        return nil, err
    }
    if !sessionExists {
        return nil, fmt.Errorf("session not found")
    }

    // Verify word exists
    var wordExists bool
    err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM words WHERE id = ?)", wordID).Scan(&wordExists)
    if err != nil {
        log.Printf("Error checking word existence: %v", err)
        return nil, err
    }
    if !wordExists {
        return nil, fmt.Errorf("word not found")
    }

    // Create the review
    result, err := db.Exec(`
        INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
        wordID, sessionID, req.IsCorrect)
    if err != nil {
        log.Printf("Error creating word review: %v", err)
        return nil, err
    }

    reviewID, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error getting last insert ID: %v", err)
        return nil, err
    }

    // Get the created review
    var review CreateWordReviewResponse
    err = db.QueryRow(`
        SELECT id, word_id, study_session_id, correct, created_at
        FROM word_review_items
        WHERE id = ?`,
        reviewID).Scan(
        &review.ID,
        &review.WordID,
        &review.SessionID,
        &review.IsCorrect,
        &review.ReviewedAt,
    )
    if err != nil {
        log.Printf("Error getting created review: %v", err)
        return nil, err
    }

    return &review, nil
}
