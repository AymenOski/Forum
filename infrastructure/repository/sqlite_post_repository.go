package repository

import (
	"database/sql"

	"forum/domain/entity"

	"github.com/google/uuid"
)

type sqlitePostRepo struct {
	db *sql.DB
}

func (r *sqlitePostRepo) GetByUserID(userID *uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT p.post_id, p.user_.id, p.author_name ,p.content
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
			&post.ID,
			&post.UserID,
			&post.AuthorName,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *sqlitePostRepo) GetByCategory(categoryID int) ([]*entity.Post, error) {
	query := `SELECT p.post_id, p.author_name, p.content
			FROM posts p
			JOIN post_categories pc ON pc.post_id = p.post_id
			WHERE pc.category_id = ?`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	var posts []*entity.Post
	for rows.Next() {
		post := &entity.Post{}
		err := rows.Scan(
			&post.ID,
			&post.AuthorName,
			&post.Content,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
