package service

import (
    "fmt"
    "log"
    "github.com/minhalzubairi/lang-portal/backend-go/internal/models"
)

// GetWords returns a paginated list of words with their stats
func GetWords(page, perPage int) ([]WordWithStats, *models.Pagination, error) {
    db := GetDB()
    if db == nil {
        log.Printf("Database connection is nil")
        return nil, nil, fmt.Errorf("database connection not initialized")
    }

    offset := (page - 1) * perPage

    // Get total count for pagination
    var total int
    err := db.QueryRow("SELECT COUNT(*) FROM words").Scan(&total)
    if err != nil {
        log.Printf("Error counting words: %v", err)
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
        ORDER BY w.id
        LIMIT ? OFFSET ?`,
        perPage, offset)
    if err != nil {
        log.Printf("Error querying words: %v", err)
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

// WordDetailResponse represents a single word with stats and groups
type WordDetailResponse struct {
    Arabic  string   `json:"arabic"`
    Roman   string   `json:"roman"`
    English string   `json:"english"`
    Stats   struct {
        CorrectCount int `json:"correct_count"`
        WrongCount   int `json:"wrong_count"`
    } `json:"stats"`
    Groups []struct {
        Name string `json:"name"`
    } `json:"groups"`
}

// GetWord returns a single word by ID with stats and groups
func GetWord(id int64) (*WordDetailResponse, error) {
    db := GetDB()
    if db == nil {
        return nil, fmt.Errorf("database connection not initialized")
    }

    var word WordDetailResponse

    // Get word details with stats
    err := db.QueryRow(`
        SELECT 
            w.arabic,
            w.roman,
            w.english,
            COALESCE(correct.count, 0) as correct_count,
            COALESCE(wrong.count, 0) as wrong_count
        FROM words w
        LEFT JOIN (
            SELECT word_id, COUNT(*) as count
            FROM word_review_items
            WHERE correct = 1 AND word_id = ?
            GROUP BY word_id
        ) correct ON w.id = correct.word_id
        LEFT JOIN (
            SELECT word_id, COUNT(*) as count
            FROM word_review_items
            WHERE correct = 0 AND word_id = ?
            GROUP BY word_id
        ) wrong ON w.id = wrong.word_id
        WHERE w.id = ?`,
        id, id, id).Scan(
            &word.Arabic,
            &word.Roman,
            &word.English,
            &word.Stats.CorrectCount,
            &word.Stats.WrongCount,
        )
    if err != nil {
        log.Printf("Error getting word %d: %v", id, err)
        return nil, err
    }

    // Get groups for this word
    rows, err := db.Query(`
        SELECT g.name
        FROM groups g
        JOIN words_groups wg ON g.id = wg.group_id
        WHERE wg.word_id = ?
        ORDER BY g.name`,
        id)
    if err != nil {
        log.Printf("Error getting word groups: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var group struct {
            Name string `json:"name"`
        }
        if err := rows.Scan(&group.Name); err != nil {
            log.Printf("Error scanning group: %v", err)
            return nil, err
        }
        word.Groups = append(word.Groups, group)
    }

    return &word, nil
}
