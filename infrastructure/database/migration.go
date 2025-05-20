package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	createUsersTable(db)
	createPostsTable(db)

	// createCommentsTable(db)
	// createCategoriesTable(db)
	// createLikesTable(db)

}

func createUsersTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		session_token TEXT,
		session_expiry DATETIME
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Failed to create users table:", err)
	}
}

func createPostsTable(db *sql.DB) {
	query := `
        CREATE TABLE IF NOT EXISTS posts (
            user_id INTEGER NOT NULL,
            post_id INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            likes_count INTEGER DEFAULT 0,
			dislikes_count INTEGER DEFAULT 0,
            FOREIGN KEY(user_id) REFERENCES users(user_id)
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("❌ Failed to create posts table:", err)
	}
}
