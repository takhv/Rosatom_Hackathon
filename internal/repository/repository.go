package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./mvp.db")
	if err != nil {
		log.Fatal("sql.Open fatal error")
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		log.Fatal("runMigrations() error")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("db.Ping problem")
		db.Close()
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatal(files)
		log.Fatal("read migrations dir error")
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println("file is dir")
			continue
		}

		content, err := os.ReadFile("migrations/" + file.Name())
		if err != nil {
			log.Fatal("error with read mifration file")
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal("error with execute sql command")
			return err
		}

		log.Printf("✅ Applied migration: %s", file.Name())
	}

	return nil
}
