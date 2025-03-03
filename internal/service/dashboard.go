package service

import (
    "database/sql"
    "fmt"
    "log"
    "time"
)

// LastStudySessionResponse represents the last study session with stats
type LastStudySessionResponse struct {
    ID           int64      `json:"id"`
    ActivityName string     `json:"activity_name"`
    GroupName    string     `json:"group_name"`
    StartTime    time.Time  `json:"start_time"`
    EndTime      *time.Time `json:"end_time,omitempty"`
    Stats        struct {
        TotalWords   int `json:"total_words"`
        CorrectCount int `json:"correct_count"`
        WrongCount   int `json:"wrong_count"`
    } `json:"stats"`
}

// DailyStats represents statistics for a single day
type DailyStats struct {
    Date         string  `json:"date"`
    CorrectCount int     `json:"correct_count"`
    WrongCount   int     `json:"wrong_count"`
}

// TotalStats represents overall study statistics
type TotalStats struct {
    TotalWordsStudied int     `json:"total_words_studied"`
    TotalCorrect      int     `json:"total_correct"`
    TotalWrong        int     `json:"total_wrong"`
    AccuracyRate      float64 `json:"accuracy_rate"`
}

// StudyProgressResponse represents the complete study progress
type StudyProgressResponse struct {
    DailyStats []DailyStats `json:"daily_stats"`
    TotalStats TotalStats   `json:"total_stats"`
}

// QuickStatsLastSession represents the last study session summary
type QuickStatsLastSession struct {
    ActivityName string `json:"activity_name"`
    GroupName    string `json:"group_name"`
    CorrectCount int    `json:"correct_count"`
    WrongCount   int    `json:"wrong_count"`
}

// QuickStatsResponse represents the dashboard quick statistics
type QuickStatsResponse struct {
    TotalWordsAvailable    int                   `json:"total_words_available"`
    WordsStudied           int                   `json:"words_studied"`
    StudySessionsCompleted int                   `json:"study_sessions_completed"`
    LastStudySession       QuickStatsLastSession `json:"last_study_session"`
}

// GetLastStudySession returns the most recent study session with stats
func GetLastStudySession() (*LastStudySessionResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    var session LastStudySessionResponse
    err := db.QueryRow(`
        SELECT 
            ss.id,
            sa.name as activity_name,
            g.name as group_name,
            ss.created_at as start_time,
            (
                SELECT COUNT(DISTINCT wri.word_id)
                FROM word_review_items wri
                WHERE wri.study_session_id = ss.id
            ) as total_words,
            (
                SELECT COUNT(*)
                FROM word_review_items wri
                WHERE wri.study_session_id = ss.id AND wri.correct = 1
            ) as correct_count,
            (
                SELECT COUNT(*)
                FROM word_review_items wri
                WHERE wri.study_session_id = ss.id AND wri.correct = 0
            ) as wrong_count
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        JOIN groups g ON ss.group_id = g.id
        ORDER BY ss.created_at DESC
        LIMIT 1`,
    ).Scan(
        &session.ID,
        &session.ActivityName,
        &session.GroupName,
        &session.StartTime,
        &session.Stats.TotalWords,
        &session.Stats.CorrectCount,
        &session.Stats.WrongCount,
    )
    if err != nil {
        log.Printf("Error getting last study session: %v", err)
        return nil, err
    }

    return &session, nil
}

// GetStudyProgress returns the study progress over time
func GetStudyProgress() (*StudyProgressResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    // Get daily stats for the last 7 days
    rows, err := db.Query(`
        SELECT 
            DATE(created_at) as study_date,
            SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END) as wrong_count
        FROM word_review_items
        WHERE created_at >= date('now', '-7 days')
        GROUP BY DATE(created_at)
        ORDER BY study_date DESC`)
    if err != nil {
        log.Printf("Error getting daily stats: %v", err)
        return nil, err
    }
    defer rows.Close()

    var response StudyProgressResponse
    for rows.Next() {
        var stats DailyStats
        if err := rows.Scan(&stats.Date, &stats.CorrectCount, &stats.WrongCount); err != nil {
            log.Printf("Error scanning daily stats: %v", err)
            return nil, err
        }
        response.DailyStats = append(response.DailyStats, stats)
    }

    // Get total stats
    err = db.QueryRow(`
        SELECT 
            COUNT(DISTINCT word_id) as total_words,
            SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END) as total_correct,
            SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END) as total_wrong
        FROM word_review_items`).Scan(
        &response.TotalStats.TotalWordsStudied,
        &response.TotalStats.TotalCorrect,
        &response.TotalStats.TotalWrong)
    if err != nil {
        log.Printf("Error getting total stats: %v", err)
        return nil, err
    }

    // Calculate accuracy rate
    totalAttempts := response.TotalStats.TotalCorrect + response.TotalStats.TotalWrong
    if totalAttempts > 0 {
        response.TotalStats.AccuracyRate = float64(response.TotalStats.TotalCorrect) / float64(totalAttempts) * 100
    }

    return &response, nil
}

// GetQuickStats returns a quick overview of learning progress
func GetQuickStats() (*QuickStatsResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    var stats QuickStatsResponse

    // Get total words available
    err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&stats.TotalWordsAvailable)
    if err != nil {
        log.Printf("Error counting total words: %v", err)
        return nil, err
    }

    // Get number of unique words studied
    err = db.QueryRow(`
        SELECT COUNT(DISTINCT word_id) 
        FROM word_review_items`).Scan(&stats.WordsStudied)
    if err != nil {
        log.Printf("Error counting studied words: %v", err)
        return nil, err
    }

    // Get total study sessions completed
    err = db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&stats.StudySessionsCompleted)
    if err != nil {
        log.Printf("Error counting study sessions: %v", err)
        return nil, err
    }

    // Get last study session details
    err = db.QueryRow(`
        SELECT 
            sa.name as activity_name,
            g.name as group_name,
            (
                SELECT COUNT(*) 
                FROM word_review_items wri 
                WHERE wri.study_session_id = ss.id AND wri.correct = 1
            ) as correct_count,
            (
                SELECT COUNT(*) 
                FROM word_review_items wri 
                WHERE wri.study_session_id = ss.id AND wri.correct = 0
            ) as wrong_count
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        JOIN groups g ON ss.group_id = g.id
        ORDER BY ss.created_at DESC
        LIMIT 1
    `).Scan(
        &stats.LastStudySession.ActivityName,
        &stats.LastStudySession.GroupName,
        &stats.LastStudySession.CorrectCount,
        &stats.LastStudySession.WrongCount,
    )
    if err != nil && err != sql.ErrNoRows {
        log.Printf("Error getting last session stats: %v", err)
        return nil, err
    }

    return &stats, nil
}
