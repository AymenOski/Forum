package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	createUsersTable(db)
	createPostsTable(db)
	createCategoriesTable(db)
	createLikesDislikesTable(db)
	createCommentsTable(db)
	createPostCategoriesTables(db)
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
		log.Fatal("Failed to create users table:", err)
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
		log.Fatal("Failed to create posts table:", err)
	}
}

func createLikesDislikesTable(db *sql.DB) {
	query := `CREATE TABLE likes_dislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id TEXT NOT NULL,
			is_like BOOLEAN NOT NULL,
			FOREIGN KEY(post_id) REFERENCES posts(post_id),
			FOREIGN KEY(user_id) REFERENCES users(user_id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create likes_dislikes table:", err)
	}
}

func createCommentsTable(db *sql.DB) {
	query := `CREATE TABLE comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			content TEXT NOT NULL,
			FOREIGN KEY(post_id) REFERENCES posts(post_id),
			FOREIGN KEY(user_id) REFERENCES users(user_id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create comments table:", err)
	}
}

func createPostCategoriesTables(db *sql.DB) {
	query := `CREATE TABLE post_categories (
			post_id TEXT NOT NULL,
			category_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, category_id),
			FOREIGN KEY(post_id) REFERENCES posts(post_id),
			FOREIGN KEY(category_id) REFERENCES categories(category_id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create post_categories table:", err)
	}
}

func createCategoriesTable(db *sql.DB) {
	query := `CREATE TABLE categories (
			category_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create categories table:", err)
	}
}
