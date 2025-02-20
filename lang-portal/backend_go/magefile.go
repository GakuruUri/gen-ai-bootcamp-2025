//go:build mage
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/magefile/mage/sh"
)

// InitDB initializes the SQLite database
func InitDB() error {
	fmt.Println("Initializing database...")
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return err
	}
	defer db.Close()
	return nil
}

// Migrate runs all database migrations
func Migrate() error {
	fmt.Println("Running migrations...")
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return err
	}
	defer db.Close()

	files, err := filepath.Glob("db/migrations/*.sql")
	if err != nil {
		return err
	}

	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Applying migration: %s\n", file)
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("error executing %s: %v", file, err)
			}
		}
	}

	return nil
}

// Seed imports sample data into the database
func Seed() error {
	fmt.Println("Seeding database...")
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return err
	}
	defer db.Close()

	files, err := filepath.Glob("db/seeds/*.json")
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Printf("Processing seed file: %s\n", file)
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		var words []struct {
			Japanese string                 `json:"japanese"`
			Romaji   string                 `json:"romaji"`
			English  string                 `json:"english"`
			Parts    map[string]interface{} `json:"parts"`
		}

		if err := json.Unmarshal(content, &words); err != nil {
			return err
		}

		// Get group name from filename
		groupName := strings.TrimSuffix(filepath.Base(file), ".json")
		
		// Insert or get group
		var groupID int64
		err = db.QueryRow("INSERT OR IGNORE INTO groups (name) VALUES (?)", groupName)
		if err != nil {
			return err
		}
		err = db.QueryRow("SELECT id FROM groups WHERE name = ?", groupName).Scan(&groupID)
		if err != nil {
			return err
		}

		// Insert words and create word-group associations
		for _, word := range words {
			partsJSON, err := json.Marshal(word.Parts)
			if err != nil {
				return err
			}

			var wordID int64
			err = db.QueryRow(`
				INSERT INTO words (japanese, romaji, english, parts)
				VALUES (?, ?, ?, ?)
				RETURNING id
			`, word.Japanese, word.Romaji, word.English, string(partsJSON)).Scan(&wordID)
			if err != nil {
				return err
			}

			_, err = db.Exec(`
				INSERT OR IGNORE INTO word_groups (word_id, group_id)
				VALUES (?, ?)
			`, wordID, groupID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Install installs project dependencies
func Install() error {
	fmt.Println("Installing dependencies...")
	return sh.Run("go", "mod", "download")
}

// Run starts the server
func Run() error {
	fmt.Println("Starting server...")
	return sh.Run("go", "run", "cmd/server/main.go")
}

// Build builds the project
func Build() error {
	fmt.Println("Building project...")
	return sh.Run("go", "build", "-o", "server", "cmd/server/main.go")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	if err := os.Remove("server"); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.Remove("words.db"); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Setup runs the complete setup process
func Setup() error {
	if err := Clean(); err != nil {
		return err
	}
	if err := Install(); err != nil {
		return err
	}
	if err := InitDB(); err != nil {
		return err
	}
	if err := Migrate(); err != nil {
		return err
	}
	if err := Seed(); err != nil {
		return err
	}
	return nil
}
