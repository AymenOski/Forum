package infra_repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/domain/entity"

	"github.com/google/uuid"
)

type SQLitePostRepository struct {
	db *sql.DB
}

func NewSQLitePostRepository(db *sql.DB) *SQLitePostRepository {
	return &SQLitePostRepository{db: db}
}

func (r *SQLitePostRepository) Create(post *entity.Post) error {
	// Generate UUID for new Post
	post.PostID = uuid.New()

	query := `INSERT INTO posts (post_id, user_id, content, created_at) 
			  VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, post.PostID, post.UserID, post.Content, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLitePostRepository) GetByID(postID uuid.UUID) (*entity.Post, error) {
	query := `SELECT p.post_id, p.user_id, u.name as author_name, p.content, 
			 p.likes_count, p.dislikes_count, p.created_at
			 FROM posts p
			 JOIN users u ON p.user_id = u.user_id
			 WHERE p.post_id = ?`

	row := r.db.QueryRow(query, postID)

	post := &entity.Post{}
	var userIDStr string

	err := row.Scan(
		&post.PostID,
		&userIDStr,
		&post.Authorname,
		&post.Content,
		&post.LikesCount,
		&post.DislikesCount,
		&post.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post not found")
		}
		return nil, err
	}

	// Parse UUID
	post.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *SQLitePostRepository) GetByUserID(userID uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT p.post_id, p.user_id, u.name as author_name, p.content,
			 p.likes_count, p.dislikes_count, p.created_at
			 FROM posts p
			 JOIN users u ON p.user_id = u.user_id
			 WHERE p.user_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		var userIDStr string

		err := rows.Scan(
			&post.PostID,
			&userIDStr,
			&post.Authorname,
			&post.Content,
			&post.LikesCount,
			&post.DislikesCount,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse UUID
		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetByCategory(categoryIDs []uuid.UUID) ([]*entity.Post, error) {
	if len(categoryIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(categoryIDs))
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT p.post_id, p.user_id, u.name as author_name, p.content,
						p.likes_count, p.dislikes_count, p.created_at
						FROM posts p
						JOIN users u ON p.user_id = u.user_id
						JOIN post_categories pc ON pc.post_id = p.post_id
						WHERE pc.category_id IN (%s)`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		var userIDStr string

		err := rows.Scan(
			&post.PostID,
			&userIDStr,
			&post.Authorname,
			&post.Content,
			&post.LikesCount,
			&post.DislikesCount,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse UUID
		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetAll() ([]*entity.Post, error) {
	query := `SELECT p.post_id, p.user_id, u.name as author_name, p.content,
			 p.likes_count, p.dislikes_count, p.created_at
			 FROM posts p
			 JOIN users u ON p.user_id = u.user_id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		var userIDStr string

		err := rows.Scan(
			&post.PostID,
			&userIDStr,
			&post.Authorname,
			&post.Content,
			&post.LikesCount,
			&post.DislikesCount,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse UUID
		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetLikedPosts(userID uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT p.post_id, p.user_id, u.name as author_name, p.content,
			 p.likes_count, p.dislikes_count, p.created_at
			 FROM posts p
			 JOIN users u ON p.user_id = u.user_id
			 JOIN post_reactions pr ON pr.post_id = p.post_id
			 WHERE pr.user_id = ? AND pr.reaction_type = 1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		var userIDStr string

		err := rows.Scan(
			&post.PostID,
			&userIDStr,
			&post.Authorname,
			&post.Content,
			&post.LikesCount,
			&post.DislikesCount,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse UUID
		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *SQLitePostRepository) Update(post *entity.Post) error {
	query := `UPDATE posts 
			 SET content = ?, likes_count = ?, dislikes_count = ?
			 WHERE post_id = ?`

	_, err := r.db.Exec(query, post.Content, post.LikesCount, post.DislikesCount, post.PostID)
	return err
}

func (r *SQLitePostRepository) Delete(postID uuid.UUID) error {
	query := `DELETE FROM posts WHERE post_id = ?`
	_, err := r.db.Exec(query, postID)
	return err
}
