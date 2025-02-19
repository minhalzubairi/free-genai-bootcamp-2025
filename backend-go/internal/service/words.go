package service

import (
	"fmt"
	"log"
	"github.com/minhalzubairi/lang-portal/backend-go/internal/models"
)

// GetWords returns a paginated list of words
func GetWords(page, perPage int) ([]models.Word, *models.Pagination, error) {
	db := GetDB() // Get the database instance
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

	// Get words with pagination
	rows, err := db.Query(`
		SELECT id, arabic, roman, english, parts 
		FROM words 
		LIMIT ? OFFSET ?`, 
		perPage, offset)
	if err != nil {
		log.Printf("Error querying words: %v", err)
		return nil, nil, err
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var w models.Word
		err := rows.Scan(&w.ID, &w.Arabic, &w.Roman, &w.English, &w.Parts)
		if err != nil {
			log.Printf("Error scanning word: %v", err)
			return nil, nil, err
		}
		// Only include parts in JSON if it's not null
		if !w.Parts.Valid {
			w.Parts.String = ""
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
