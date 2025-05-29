package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) {
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal("Failed to enable foreign key constraints:", err)
	}
	createUsersTable(db)
	createPostsTable(db)
	createCommentsTable(db)
	createCategoriesTable(db)
	createPostCategoriesTable(db)
	createUserSessionsTable(db)
	createCommentReactionTable(db)
	createPostReactionTable(db)
}

func createUsersTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS user (
		id CHAR(36) NOT NULL,
		user_name TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY(id)
	);
	CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create user table:", err)
	}
}

func createPostsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id CHAR(36) NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id CHAR(36) NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY(id),
		FOREIGN KEY(user_id) REFERENCES user(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create posts table:", err)
	}
}

func createCommentsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS comments (
		id CHAR(36) NOT NULL,
		content TEXT NOT NULL,
		user_id CHAR(36) NOT NULL,
		post_id CHAR(36) NOT NULL,
		createdat DATETIME NOT NULL,
		PRIMARY KEY(id),
		FOREIGN KEY(user_id) REFERENCES user(id),
		FOREIGN KEY(post_id) REFERENCES posts(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create comments table:", err)
	}
}

func createCategoriesTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id CHAR(36) NOT NULL,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create categories table:", err)
	}
}

func createPostCategoriesTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS post_categories (
		post_id CHAR(36) NOT NULL,
		category_id CHAR(36) NOT NULL,
		PRIMARY KEY(post_id, category_id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create post_categories table:", err)
	}
}

func createUserSessionsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS user_sessions (
		id CHAR(36) NOT NULL,
		user_id CHAR(36) NOT NULL,
		session_token TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		created_at DATETIME NOT NULL,
		PRIMARY KEY(id),
		FOREIGN KEY(user_id) REFERENCES user(id)
	);
	CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions(session_token);
	CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create user_sessions table:", err)
	}
}

func createCommentReactionTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS comment_reaction (
		id CHAR(36) NOT NULL,
		user_id CHAR(36) NOT NULL,
		comment_id CHAR(36) NOT NULL,
		reaction INTEGER NOT NULL CHECK (reaction IN (0, 1)),
		created_at DATETIME NOT NULL,
		PRIMARY KEY(id),
		FOREIGN KEY(user_id) REFERENCES user(id),
		FOREIGN KEY(comment_id) REFERENCES comments(id),
		UNIQUE(user_id, comment_id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create comment_reaction table:", err)
	}
}

func createPostReactionTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS post_reaction (
		id CHAR(36) NOT NULL,
		user_id CHAR(36) NOT NULL,
		post_id CHAR(36) NOT NULL,
		reaction INTEGER NOT NULL CHECK (reaction IN (0, 1)),
		created_at DATETIME NOT NULL,
		PRIMARY KEY(id),
		FOREIGN KEY(user_id) REFERENCES user(id),
		FOREIGN KEY(post_id) REFERENCES posts(id),
		UNIQUE(user_id, post_id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create post_reaction table:", err)
	}
}