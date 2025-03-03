package service

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    "github.com/minhalzubairi/lang-portal/backend-go/internal/models"
)

// ActivityStats represents statistics for a study activity
type ActivityStats struct {
    TotalSessions      int     `json:"total_sessions"`
    TotalWordsReviewed int     `json:"total_words_reviewed"`
    AccuracyRate       float64 `json:"accuracy_rate"`
}

// ActivityResponse represents a study activity with its stats
type ActivityResponse struct {
    ID           int64         `json:"id"`
    Name         string        `json:"name"`
    ThumbnailURL string        `json:"thumbnail_url"`
    Description  string        `json:"description"`
    LaunchURL    string        `json:"launch_url"`
    Stats        ActivityStats `json:"stats"`
}

// ActivitySessionStats represents statistics for a study session
type ActivitySessionStats struct {
    TotalWords   int `json:"total_words"`
    CorrectCount int `json:"correct_count"`
    WrongCount   int `json:"wrong_count"`
}

// RecentSession represents a recent study session
type RecentSession struct {
    ID        int64              `json:"id"`
    GroupName string             `json:"group_name"`
    StartTime time.Time          `json:"start_time"`
    Stats     ActivitySessionStats `json:"stats"`
}

// ActivityDetailResponse represents a single study activity with detailed stats
type ActivityDetailResponse struct {
    ID              int64           `json:"id"`
    Name            string          `json:"name"`
    ThumbnailURL    string          `json:"thumbnail_url"`
    Description     string          `json:"description"`
    LaunchURL       string          `json:"launch_url"`
    Stats           ActivityStats   `json:"stats"`
    RecentSessions []RecentSession `json:"recent_sessions"`
}

// GetStudyActivities returns a paginated list of study activities with their stats
func GetStudyActivities(page, perPage int) ([]ActivityResponse, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow("SELECT COUNT(*) FROM study_activities").Scan(&total)
    if err != nil {
        log.Printf("Error counting activities: %v", err)
        return nil, nil, err
    }

    // Get activities with their stats
    rows, err := db.Query(`
        SELECT 
            sa.id,
            sa.name,
            sa.thumbnail_url,
            sa.description,
            sa.launch_url,
            COUNT(DISTINCT ss.id) as total_sessions,
            COUNT(wri.id) as total_reviews,
            COALESCE(
                SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) * 100.0 / 
                NULLIF(COUNT(wri.id), 0),
                0
            ) as accuracy_rate
        FROM study_activities sa
        LEFT JOIN study_sessions ss ON sa.id = ss.study_activity_id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        GROUP BY sa.id
        ORDER BY sa.name
        LIMIT ? OFFSET ?`,
        perPage, offset)
    if err != nil {
        log.Printf("Error querying activities: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var activities []ActivityResponse
    for rows.Next() {
        var a ActivityResponse
        err := rows.Scan(
            &a.ID,
            &a.Name,
            &a.ThumbnailURL,
            &a.Description,
            &a.LaunchURL,
            &a.Stats.TotalSessions,
            &a.Stats.TotalWordsReviewed,
            &a.Stats.AccuracyRate,
        )
        if err != nil {
            log.Printf("Error scanning activity: %v", err)
            return nil, nil, err
        }
        activities = append(activities, a)
    }

    pagination := &models.Pagination{
        CurrentPage:  page,
        ItemsPerPage: perPage,
        TotalItems:   total,
        TotalPages:   (total + perPage - 1) / perPage,
    }

    return activities, pagination, nil
}

// GetStudyActivity returns a single study activity by ID with detailed stats
func GetStudyActivity(id int64) (*ActivityDetailResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    var activity ActivityDetailResponse

    // Get activity details with overall stats
    err := db.QueryRow(`
        SELECT 
            sa.id,
            sa.name,
            sa.thumbnail_url,
            sa.description,
            sa.launch_url,
            COUNT(DISTINCT ss.id) as total_sessions,
            COUNT(wri.id) as total_reviews,
            COALESCE(
                SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) * 100.0 / 
                NULLIF(COUNT(wri.id), 0),
                0
            ) as accuracy_rate
        FROM study_activities sa
        LEFT JOIN study_sessions ss ON sa.id = ss.study_activity_id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        WHERE sa.id = ?
        GROUP BY sa.id`,
        id).Scan(
            &activity.ID,
            &activity.Name,
            &activity.ThumbnailURL,
            &activity.Description,
            &activity.LaunchURL,
            &activity.Stats.TotalSessions,
            &activity.Stats.TotalWordsReviewed,
            &activity.Stats.AccuracyRate,
        )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, err
        }
        log.Printf("Error getting activity %d: %v", id, err)
        return nil, err
    }

    // Get recent sessions (last 5)
    rows, err := db.Query(`
        SELECT 
            ss.id,
            g.name as group_name,
            ss.created_at as start_time,
            COUNT(DISTINCT wri.word_id) as total_words,
            SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END) as wrong_count
        FROM study_sessions ss
        JOIN groups g ON ss.group_id = g.id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        WHERE ss.study_activity_id = ?
        GROUP BY ss.id
        ORDER BY ss.created_at DESC
        LIMIT 5`,
        id)
    if err != nil {
        log.Printf("Error getting recent sessions: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var session RecentSession
        err := rows.Scan(
            &session.ID,
            &session.GroupName,
            &session.StartTime,
            &session.Stats.TotalWords,
            &session.Stats.CorrectCount,
            &session.Stats.WrongCount,
        )
        if err != nil {
            log.Printf("Error scanning session: %v", err)
            return nil, err
        }
        activity.RecentSessions = append(activity.RecentSessions, session)
    }

    return &activity, nil
}

// ActivitySessionResponse represents a study session for an activity
type ActivitySessionResponse struct {
    ID        int64      `json:"id"`
    GroupName string     `json:"group_name"`
    StartTime time.Time  `json:"start_time"`
    EndTime   *time.Time `json:"end_time,omitempty"`
    Stats     struct {
        TotalWords   int `json:"total_words"`
        CorrectCount int `json:"correct_count"`
        WrongCount   int `json:"wrong_count"`
    } `json:"stats"`
}

// GetStudyActivitySessions returns paginated study sessions for an activity
func GetStudyActivitySessions(activityID int64, page, perPage int) ([]ActivitySessionResponse, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow(`
        SELECT COUNT(*)
        FROM study_sessions
        WHERE study_activity_id = ?`,
        activityID).Scan(&total)
    if err != nil {
        log.Printf("Error counting activity sessions: %v", err)
        return nil, nil, err
    }

    // Get sessions with stats
    rows, err := db.Query(`
        SELECT 
            ss.id,
            g.name as group_name,
            ss.created_at as start_time,
            COUNT(DISTINCT wri.word_id) as total_words,
            SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END) as correct_count,
            SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END) as wrong_count
        FROM study_sessions ss
        JOIN groups g ON ss.group_id = g.id
        LEFT JOIN word_review_items wri ON ss.id = wri.study_session_id
        WHERE ss.study_activity_id = ?
        GROUP BY ss.id
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?`,
        activityID, perPage, offset)
    if err != nil {
        log.Printf("Error querying activity sessions: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var sessions []ActivitySessionResponse
    for rows.Next() {
        var s ActivitySessionResponse
        err := rows.Scan(
            &s.ID,
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

// CreateActivityRequest represents the request body for creating a new activity
type CreateActivityRequest struct {
    Name         string `json:"name" binding:"required"`
    ThumbnailURL string `json:"thumbnail_url" binding:"required"`
    Description  string `json:"description" binding:"required"`
    LaunchURL    string `json:"launch_url" binding:"required"`
}

// CreateActivityResponse represents the response for a newly created activity
type CreateActivityResponse struct {
    ID           int64  `json:"id"`
    Name         string `json:"name"`
    ThumbnailURL string `json:"thumbnail_url"`
    Description  string `json:"description"`
    LaunchURL    string `json:"launch_url"`
}

// CreateStudyActivity creates a new study activity
func CreateStudyActivity(req *CreateActivityRequest) (*CreateActivityResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    result, err := db.Exec(`
        INSERT INTO study_activities (name, thumbnail_url, description, launch_url)
        VALUES (?, ?, ?, ?)`,
        req.Name, req.ThumbnailURL, req.Description, req.LaunchURL)
    if err != nil {
        log.Printf("Error creating activity: %v", err)
        return nil, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        log.Printf("Error getting last insert ID: %v", err)
        return nil, err
    }

    return &CreateActivityResponse{
        ID:           id,
        Name:         req.Name,
        ThumbnailURL: req.ThumbnailURL,
        Description:  req.Description,
        LaunchURL:    req.LaunchURL,
    }, nil
}
