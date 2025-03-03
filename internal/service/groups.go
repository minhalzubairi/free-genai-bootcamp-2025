package service

import (
    "fmt"
    "log"
    "time"
    "github.com/minhalzubairi/lang-portal/backend-go/internal/models"
)

// GroupResponse represents a group with word count
type GroupResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    WordCount int    `json:"word_count"`
}

// GroupDetailResponse represents a single group with stats
type GroupDetailResponse struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Stats struct {
        TotalWordCount int `json:"total_word_count"`
    } `json:"stats"`
}

// WordWithStats represents a word with its review statistics
type WordWithStats struct {
    Arabic       string `json:"arabic"`
    Roman        string `json:"roman"`
    English      string `json:"english"`
    CorrectCount int    `json:"correct_count"`
    WrongCount   int    `json:"wrong_count"`
}

// GetGroups returns all word groups
func GetGroups(page, perPage int) ([]GroupResponse, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&total)
    if err != nil {
        log.Printf("Error counting groups: %v", err)
        return nil, nil, err
    }

    // Get groups with word count
    rows, err := db.Query(`
        SELECT 
            g.id,
            g.name,
            (SELECT COUNT(*) FROM words_groups wg WHERE wg.group_id = g.id) as word_count
        FROM groups g
        ORDER BY g.name
        LIMIT ? OFFSET ?`,
        perPage, offset)
    if err != nil {
        log.Printf("Error querying groups: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var groups []GroupResponse
    for rows.Next() {
        var g GroupResponse
        if err := rows.Scan(&g.ID, &g.Name, &g.WordCount); err != nil {
            log.Printf("Error scanning group: %v", err)
            return nil, nil, err
        }
        groups = append(groups, g)
    }

    pagination := &models.Pagination{
        CurrentPage:  page,
        ItemsPerPage: perPage,
        TotalItems:   total,
        TotalPages:   (total + perPage - 1) / perPage,
    }

    return groups, pagination, nil
}

// GetGroup returns a single group by ID with stats
func GetGroup(id int64) (*GroupDetailResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    var group GroupDetailResponse
    err := db.QueryRow(`
        SELECT 
            g.id,
            g.name,
            (SELECT COUNT(*) FROM words_groups wg WHERE wg.group_id = g.id) as word_count
        FROM groups g
        WHERE g.id = ?`, id).Scan(&group.ID, &group.Name, &group.Stats.TotalWordCount)
    if err != nil {
        log.Printf("Error getting group %d: %v", id, err)
        return nil, err
    }

    return &group, nil
}

// GetGroupWords returns paginated words in a group with their stats
func GetGroupWords(groupID int64, page, perPage int) ([]WordWithStats, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow(`
        SELECT COUNT(*)
        FROM words_groups
        WHERE group_id = ?`, groupID).Scan(&total)
    if err != nil {
        log.Printf("Error counting group words: %v", err)
        return nil, nil, err
    }

    // Get words with their stats
    rows, err := db.Query(`
        SELECT 
            w.arabic,
            w.roman,
            w.english,
            COALESCE(correct.count, 0) as correct_count,
            COALESCE(wrong.count, 0) as wrong_count
        FROM words w
        JOIN words_groups wg ON w.id = wg.word_id
        LEFT JOIN (
            SELECT word_id, COUNT(*) as count
            FROM word_review_items
            WHERE correct = 1
            GROUP BY word_id
        ) correct ON w.id = correct.word_id
        LEFT JOIN (
            SELECT word_id, COUNT(*) as count
            FROM word_review_items
            WHERE correct = 0
            GROUP BY word_id
        ) wrong ON w.id = wrong.word_id
        WHERE wg.group_id = ?
        ORDER BY w.id
        LIMIT ? OFFSET ?`,
        groupID, perPage, offset)
    if err != nil {
        log.Printf("Error querying group words: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var words []WordWithStats
    for rows.Next() {
        var w WordWithStats
        if err := rows.Scan(
            &w.Arabic,
            &w.Roman,
            &w.English,
            &w.CorrectCount,
            &w.WrongCount,
        ); err != nil {
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

// GroupStudySession represents a study session with activity details
type GroupStudySession struct {
    ID               int64     `json:"id"`
    ActivityName     string    `json:"activity_name"`
    GroupName        string    `json:"group_name"`
    StartTime        time.Time `json:"start_time"`
    EndTime         *time.Time `json:"end_time,omitempty"`
    ReviewItemsCount int       `json:"review_items_count"`
}

// GetGroupStudySessions returns paginated study sessions for a group
func GetGroupStudySessions(groupID int64, page, perPage int) ([]GroupStudySession, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count
    var total int
    err := db.QueryRow(`
        SELECT COUNT(*)
        FROM study_sessions ss
        WHERE ss.group_id = ?`, groupID).Scan(&total)
    if err != nil {
        log.Printf("Error counting group study sessions: %v", err)
        return nil, nil, err
    }

    // Get study sessions with related data
    rows, err := db.Query(`
        SELECT 
            ss.id,
            sa.name as activity_name,
            g.name as group_name,
            ss.created_at as start_time,
            (SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_session_id = ss.id) as review_count
        FROM study_sessions ss
        JOIN study_activities sa ON ss.study_activity_id = sa.id
        JOIN groups g ON ss.group_id = g.id
        WHERE ss.group_id = ?
        ORDER BY ss.created_at DESC
        LIMIT ? OFFSET ?`,
        groupID, perPage, offset)
    if err != nil {
        log.Printf("Error querying group study sessions: %v", err)
        return nil, nil, err
    }
    defer rows.Close()

    var sessions []GroupStudySession
    for rows.Next() {
        var s GroupStudySession
        if err := rows.Scan(
            &s.ID,
            &s.ActivityName,
            &s.GroupName,
            &s.StartTime,
            &s.ReviewItemsCount,
        ); err != nil {
            log.Printf("Error scanning study session: %v", err)
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
