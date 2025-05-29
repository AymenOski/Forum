package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLitePostRepository implements PostRepository interface
type SQLitePostRepository struct {
	db *sql.DB
}

func NewSQLitePostRepository(db *sql.DB) repository.PostRepository {
	return &SQLitePostRepository{db: db}
}

func (r *SQLitePostRepository) Create(post *entity.Post) error {
	post.ID = uuid.New()
	post.CreatedAt = time.Now()

	query := `INSERT INTO posts (id, title, content, user_id, created_at)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, post.ID.String(), post.Title, post.Content,
		post.UserID.String(), post.CreatedAt)
	return err
}

func (r *SQLitePostRepository) GetByID(postID uuid.UUID) (*entity.Post, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts WHERE id = ?`

	row := r.db.QueryRow(query, postID.String())

	post := &entity.Post{}
	var idStr, userIDStr string

	err := row.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	post.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	post.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (r *SQLitePostRepository) GetAll() ([]*entity.Post, error) {
	query := `SELECT id, title, content, user_id, created_at FROM posts ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetByUserID(userID uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT id, title, content, user_id, created_at 
			  FROM posts WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetWithPagination(limit, offset int) ([]*entity.Post, error) {
	query := `SELECT id, title, content, user_id, created_at 
			  FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetByCategory(categoryID uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at 
			  FROM posts p 
			  INNER JOIN post_categories pc ON p.id = pc.post_id 
			  WHERE pc.category_id = ? 
			  ORDER BY p.created_at DESC`

	rows, err := r.db.Query(query, categoryID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetByCategoryWithPagination(categoryID uuid.UUID, limit, offset int) ([]*entity.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at 
			  FROM posts p 
			  INNER JOIN post_categories pc ON p.id = pc.post_id 
			  WHERE pc.category_id = ? 
			  ORDER BY p.created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, categoryID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetMostLiked(limit int) ([]*entity.Post, error) {
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at 
			  FROM posts p 
			  LEFT JOIN (
				  SELECT post_id, COUNT(*) as like_count 
				  FROM post_reaction 
				  WHERE reaction = 1 
				  GROUP BY post_id
			  ) lr ON p.id = lr.post_id 
			  ORDER BY COALESCE(lr.like_count, 0) DESC, p.created_at DESC 
			  LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetRecent(limit int) ([]*entity.Post, error) {
	query := `SELECT id, title, content, user_id, created_at 
			  FROM posts ORDER BY created_at DESC LIMIT ?`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) Update(post *entity.Post) error {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`

	_, err := r.db.Exec(query, post.Title, post.Content, post.ID.String())
	return err
}

func (r *SQLitePostRepository) Delete(postID uuid.UUID) error {
	query := `DELETE FROM posts WHERE id = ?`

	_, err := r.db.Exec(query, postID.String())
	return err
}

func (r *SQLitePostRepository) Search(query string) ([]*entity.Post, error) {
	searchQuery := `SELECT id, title, content, user_id, created_at 
					FROM posts 
					WHERE title LIKE ? OR content LIKE ? 
					ORDER BY created_at DESC`

	searchTerm := "%" + query + "%"
	rows, err := r.db.Query(searchQuery, searchTerm, searchTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.Post

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			return nil, err
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *SQLitePostRepository) GetWithDetails(postID uuid.UUID) (*entity.PostWithDetails, error) {
	// Basic implementation - can be extended to include more details
	post, err := r.GetByID(postID)
	if err != nil {
		return nil, err
	}

	return &entity.PostWithDetails{
		Post: *post,
	}, nil
}

func (r *SQLitePostRepository) GetAllWithDetails() ([]*entity.PostWithDetails, error) {
	// Basic implementation - can be extended to include more details
	posts, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var postsWithDetails []*entity.PostWithDetails
	for _, post := range posts {
		postsWithDetails = append(postsWithDetails, &entity.PostWithDetails{
			Post: *post,
		})
	}

	return postsWithDetails, nil
}
