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

func SeedTestData(db *sql.DB) {
	seedUsers(db)
	seedCategories(db)
	seedPosts(db)
	seedLikesDislikes(db)
	seedComments(db)
	seedPostCategories(db)
}

func createUsersTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    session_token TEXT UNIQUE,
    session_expiry DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_session_token ON users(session_token);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
}

func createPostsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS posts (
		user_id TEXT NOT NULL,
		post_id INTEGER PRIMARY KEY AUTOINCREMENT,
		author_name TEXT NOT NULL,
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
	query := `CREATE TABLE reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER NOT NULL,
    user_id TEXT NOT NULL,
    is_like BOOLEAN NOT NULL,
    is_post BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(parent_id, user_id, is_post)
);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create reactions table:", err)
	}
}

func createCommentsTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS comments (
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
	query := `CREATE TABLE IF NOT EXISTS post_categories (
		post_id TEXT NOT NULL,
		category_id INTEGER NOT NULL,
		FOREIGN KEY(post_id) REFERENCES posts(post_id),
		FOREIGN KEY(category_id) REFERENCES categories(category_id)
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create post_categories table:", err)
	}
}

func createCategoriesTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS categories (
		category_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create categories table:", err)
	}
}

// Seed functions for test data

func seedUsers(db *sql.DB) {
	users := []struct {
		userID        string
		name          string
		email         string
		passwordHash  string
		sessionToken  *string
		sessionExpiry *string
	}{
		// Users with no active session
		{"550e8400-e29b-41d4-a716-446655440000", "John Doe", "john.doe@email.com", "hash1", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440001", "Bob Johnson", "bob.johnson@email.com", "hash3", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440002", "Charlie Wilson", "charlie.wilson@email.com", "hash5", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440003", "Frank Miller", "frank.miller@email.com", "hash7", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440004", "Henry Moore", "henry.moore@email.com", "hash9", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440005", "Jack Black", "jack.black@email.com", "hash11", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440006", "Liam Blue", "liam.blue@email.com", "hash13", nil, nil},
		{"550e8400-e29b-41d4-a716-446655440007", "Noah Gray", "noah.gray@email.com", "hash15", nil, nil},

		// Users with active sessions (expiring at end of 2025)
		{"550e8400-e29b-41d4-a716-446655440008", "Jane Smith", "jane.smith@email.com", "hash2", stringPtr("token2"), stringPtr("2025-12-31T23:59:59Z")},
		{"550e8400-e29b-41d4-a716-446655440009", "Alice Brown", "alice.brown@email.com", "hash4", stringPtr("token4"), stringPtr("2025-12-31T23:59:59Z")},
		{"550e8400-e29b-41d4-a716-446655440010", "Diana Davis", "diana.davis@email.com", "hash6", stringPtr("token6"), stringPtr("2025-12-31T23:59:59Z")},
		{"550e8400-e29b-41d4-a716-446655440011", "Grace Taylor", "grace.taylor@email.com", "hash8", stringPtr("token8"), stringPtr("2025-12-31T23:59:59Z")},
		{"550e8400-e29b-41d4-a716-446655440012", "Ivy White", "ivy.white@email.com", "hash10", stringPtr("token10"), stringPtr("2025-12-31T23:59:59Z")},
		{"550e8400-e29b-41d4-a716-446655440013", "Kelly Green", "kelly.green@email.com", "hash12", stringPtr("token12"), stringPtr("2025-12-31T23:59:59Z")},
		{"550e8400-e29b-41d4-a716-446655440014", "Mia Red", "mia.red@email.com", "hash14", stringPtr("token14"), stringPtr("2025-12-31T23:59:59Z")},
	}

	for _, user := range users {
		query := `INSERT OR IGNORE INTO users 
                  (user_id, name, email, password_hash, session_token, session_expiry)
                  VALUES (?, ?, ?, ?, ?, ?)`
		_, err := db.Exec(query,
			user.userID,
			user.name,
			user.email,
			user.passwordHash,
			user.sessionToken,
			user.sessionExpiry,
		)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user.name, err)
		}
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

func seedCategories(db *sql.DB) {
	categories := []string{
		"Technology", "Science", "Sports", "Entertainment", "Politics",
		"Health", "Travel", "Food", "Fashion", "Music",
		"Art", "Education", "Business", "Gaming", "Photography",
	}

	for _, category := range categories {
		query := `INSERT OR IGNORE INTO categories (name) VALUES (?)`
		_, err := db.Exec(query, category)
		if err != nil {
			log.Printf("Failed to insert category %s: %v", category, err)
		}
	}
}

func seedPosts(db *sql.DB) {
	posts := []struct {
		userID        string
		content       string
		likesCount    int
		dislikesCount int
	}{
		{"user1", "Just learned about Go programming language. It's amazing!", 5, 1},
		{"user2", "The latest scientific discovery about black holes is fascinating!", 8, 0},
		{"user3", "What a game last night! The championship was incredible.", 12, 2},
		{"user4", "New movie releases this month are looking promising.", 6, 1},
		{"user5", "The recent political developments are quite concerning.", 3, 7},
		{"user6", "Healthy eating habits can change your life completely.", 15, 0},
		{"user7", "Traveling to Japan was the best experience of my life!", 20, 1},
		{"user8", "This recipe for chocolate cake is absolutely divine.", 9, 0},
		{"user9", "Fashion trends for this season are quite interesting.", 4, 2},
		{"user10", "The concert last night was beyond expectations!", 11, 0},
		{"user11", "Modern art installations are becoming more creative.", 7, 3},
		{"user12", "Online education is revolutionizing how we learn.", 13, 1},
		{"user13", "Starting a small business requires careful planning.", 8, 0},
		{"user14", "The new gaming console has amazing graphics!", 16, 2},
		{"user15", "Photography tips for beginners: lighting is everything!", 10, 0},
	}

	for _, post := range posts {
		query := `INSERT INTO posts (user_id, content, likes_count, dislikes_count) 
				  VALUES (?, ?, ?, ?)`
		_, err := db.Exec(query, post.userID, post.content, post.likesCount, post.dislikesCount)
		if err != nil {
			log.Printf("Failed to insert post: %v", err)
		}
	}
}

func seedLikesDislikes(db *sql.DB) {
	likesData := []struct {
		postID int
		userID string
		isLike bool
	}{
		{1, "user2", true},
		{1, "user3", true},
		{1, "user4", true},
		{1, "user5", true},
		{1, "user6", true},
		{2, "user1", true},
		{2, "user3", true},
		{2, "user4", true},
		{2, "user5", true},
		{2, "user6", true},
		{3, "user1", true},
		{3, "user2", true},
		{3, "user4", true},
		{3, "user5", false},
		{3, "user6", true},
		{4, "user1", true},
		{4, "user2", true},
		{4, "user3", true},
		{4, "user5", false},
		{4, "user6", true},
		{5, "user1", false},
		{5, "user2", false},
		{5, "user3", true},
		{5, "user4", false},
		{5, "user6", false},
		{6, "user1", true},
		{6, "user2", true},
		{6, "user3", true},
		{6, "user4", true},
		{6, "user5", true},
		{7, "user1", true},
		{7, "user2", true},
		{7, "user3", true},
		{7, "user4", true},
		{7, "user5", true},
		{8, "user1", true},
		{8, "user2", true},
		{8, "user3", true},
		{8, "user4", true},
		{8, "user5", true},
		{9, "user1", true},
		{9, "user2", false},
		{9, "user3", true},
		{9, "user4", false},
		{9, "user5", true},
		{10, "user1", true},
		{10, "user2", true},
		{10, "user3", true},
		{10, "user4", true},
		{10, "user5", true},
	}

	for _, like := range likesData {
		query := `INSERT INTO likes_dislikes (post_id, user_id, is_like) VALUES (?, ?, ?)`
		_, err := db.Exec(query, like.postID, like.userID, like.isLike)
		if err != nil {
			log.Printf("Failed to insert like/dislike: %v", err)
		}
	}
}

func seedComments(db *sql.DB) {
	comments := []struct {
		postID  string
		userID  string
		content string
	}{
		{"1", "user2", "I totally agree! Go is very efficient."},
		{"1", "user3", "Have you tried building web servers with it?"},
		{"1", "user4", "The concurrency model is what makes it special."},
		{"2", "user1", "Science never fails to amaze me!"},
		{"2", "user3", "The images from the telescope are breathtaking."},
		{"3", "user1", "That final play was unbelievable!"},
		{"3", "user2", "Best championship game in years."},
		{"4", "user1", "I'm definitely watching the new thriller."},
		{"4", "user3", "The trailer looked amazing!"},
		{"5", "user1", "We need more transparency in government."},
		{"6", "user2", "Thanks for sharing these tips!"},
		{"6", "user3", "I've been following this diet for months."},
		{"7", "user1", "Japan is on my bucket list too!"},
		{"7", "user2", "The culture there is so rich and beautiful."},
		{"8", "user1", "I'm definitely trying this recipe tonight!"},
	}

	for _, comment := range comments {
		query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
		_, err := db.Exec(query, comment.postID, comment.userID, comment.content)
		if err != nil {
			log.Printf("Failed to insert comment: %v", err)
		}
	}
}

func seedPostCategories(db *sql.DB) {
	postCategories := []struct {
		postID     string
		categoryID int
	}{
		{"1", 1},   // Technology
		{"2", 2},   // Science
		{"3", 3},   // Sports
		{"4", 4},   // Entertainment
		{"5", 5},   // Politics
		{"6", 6},   // Health
		{"7", 7},   // Travel
		{"8", 8},   // Food
		{"9", 9},   // Fashion
		{"10", 10}, // Music
		{"11", 11}, // Art
		{"12", 12}, // Education
		{"13", 13}, // Business
		{"14", 14}, // Gaming
		{"15", 15}, // Photography
	}

	for _, pc := range postCategories {
		query := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`
		_, err := db.Exec(query, pc.postID, pc.categoryID)
		if err != nil {
			log.Printf("Failed to insert post category: %v", err)
		}
	}
}
