package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func SetingUpDB(datapath string) *sql.DB {
	db, err := OpenDB(datapath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	RunMigrations(db)
	return db
}

func OpenDB(filePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filePath)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseDB(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Println("Error closing database:", err)
	}
}
