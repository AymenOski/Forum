package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type sqlitePostRepo struct {
	db *sql.DB
}

func NewSqlitePostRepository(db *sql.DB) repository.PostRepository {
	return &sqlitePostRepo{
		db: db,
	}
}

func (r *sqlitePostRepo) Create(post *entity.Post) error {
	return nil
}

func (r *sqlitePostRepo) GetAll() ([]*entity.Post, error) {
	return nil, nil
}

func (r *sqlitePostRepo) GetLikedByUser(userID *uuid.UUID) ([]*entity.Post, error) {
return nil, nil
}

func (r *sqlitePostRepo) GetByUserID(userID *uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT p.post_id, p.user_id, u.name as author_name, p.content
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
		err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Authorname,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *sqlitePostRepo) GetByCategory(categoryIDs []uint8) ([]*entity.Post, error) {
	holders := make([]string, len(categoryIDs))
	args := make([]uint8, len(categoryIDs))
	for i, id := range categoryIDs {
		holders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT p.post_id, p.author_name, p.content
			FROM posts p
			JOIN post_categories pc ON pc.post_id = p.post_id
			WHERE pc.category_id IN (%s)`, strings.Join(holders, ","))

	rows, err := r.db.Query(query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		err := rows.Scan(
			&post.PostID,
			&post.Authorname,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
