//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "words.db"

// InitDB initializes the SQLite database
func InitDB() error {
	if _, err := os.Stat(dbName); err == nil {
		fmt.Println("Database already exists")
		return nil
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}
	defer db.Close()

	fmt.Println("Database initialized successfully")
	return nil
}

// Migrate runs all pending migrations
func Migrate() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create migrations table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}
	sort.Strings(files)

	// Apply each migration
	for _, file := range files {
		name := filepath.Base(file)
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM migrations WHERE name = ?)", name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if exists {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", name, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Split the file content into individual statements
		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %w", name, err)
			}
		}

		if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", name, err)
		}

		fmt.Printf("Applied migration: %s\n", name)
	}

	return nil
}

// Seed imports seed data into the database
func Seed() error {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create Basic Greetings group
	result, err := db.Exec("INSERT INTO groups (name) VALUES (?)", "Basic Greetings")
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get group id: %w", err)
	}

	// Read and parse seed file
	data, err := os.ReadFile("db/seeds/basic_greetings.json")
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	var words []struct {
		Arabic  string `json:"arabic"`
		Roman   string `json:"roman"`
		English string `json:"english"`
	}

	if err := json.Unmarshal(data, &words); err != nil {
		return fmt.Errorf("failed to parse seed file: %w", err)
	}

	// Insert words and create word-group associations
	for _, word := range words {
		result, err := db.Exec(
			"INSERT INTO words (arabic, roman, english) VALUES (?, ?, ?)",
			word.Arabic, word.Roman, word.English,
		)
		if err != nil {
			return fmt.Errorf("failed to insert word: %w", err)
		}

		wordID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get word id: %w", err)
		}

		_, err = db.Exec(
			"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
			wordID, groupID,
		)
		if err != nil {
			return fmt.Errorf("failed to create word-group association: %w", err)
		}
	}

	fmt.Println("Seed data imported successfully")
	return nil
}
