package infra_repository

import (
	"database/sql"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)
type SQLitePostCategoryRepository struct {
	db *sql.DB
}

func NewSQLitePostCategoryRepository(db *sql.DB) repository.PostCategoryRepository {
	return &SQLitePostCategoryRepository{db: db}
}

func (r *SQLitePostCategoryRepository) Create(postCategory *entity.PostCategory) error {
	query := `INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`

	_, err := r.db.Exec(query, postCategory.PostID.String(), postCategory.CategoryID.String())
	return err
}

func (r *SQLitePostCategoryRepository) Delete(postID, categoryID uuid.UUID) error {
	query := `DELETE FROM post_categories WHERE post_id = ? AND category_id = ?`

	_, err := r.db.Exec(query, postID.String(), categoryID.String())
	return err
}

func (r *SQLitePostCategoryRepository) DeleteByPostID(postID uuid.UUID) error {
	query := `DELETE FROM post_categories WHERE post_id = ?`

	_, err := r.db.Exec(query, postID.String())
	return err
}

func (r *SQLitePostCategoryRepository) DeleteByCategoryID(categoryID uuid.UUID) error {
	query := `DELETE FROM post_categories WHERE category_id = ?`

	_, err := r.db.Exec(query, categoryID.String())
	return err
}



func (r *SQLitePostCategoryRepository) GetCategoriesByPostID(postID uuid.UUID) ([]*entity.Category, error) {
	query := `SELECT c.id, c.name, c.created_at 
			  FROM categories c 
			  INNER JOIN post_categories pc ON c.id = pc.category_id 
			  WHERE pc.post_id = ? 
			  ORDER BY c.name ASC`

	rows, err := r.db.Query(query, postID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category

	for rows.Next() {
		category := &entity.Category{}
		var idStr string

		err := rows.Scan(&idStr, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, err
		}

		category.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *SQLitePostCategoryRepository) GetPostsByCategoryID(categoryID uuid.UUID) ([]*entity.Post, error) {
	query := `SELECT p.id, p.content, p.user_id, p.created_at 
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

		err := rows.Scan(&idStr, &post.Content, &userIDStr, &post.CreatedAt)
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

func (r *SQLitePostCategoryRepository) CheckAssociationExists(postID, categoryID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM post_categories WHERE post_id = ? AND category_id = ?`

	var count int
	err := r.db.QueryRow(query, postID.String(), categoryID.String()).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *SQLitePostCategoryRepository) GetAllAssociations() ([]*entity.PostCategory, error) {
	query := `SELECT post_id, category_id FROM post_categories`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postCategories []*entity.PostCategory

	for rows.Next() {
		postCategory := &entity.PostCategory{}
		var postIDStr, categoryIDStr string

		err := rows.Scan(&postIDStr, &categoryIDStr)
		if err != nil {
			return nil, err
		}

		postCategory.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}

		postCategory.CategoryID, err = uuid.Parse(categoryIDStr)
		if err != nil {
			return nil, err
		}

		postCategories = append(postCategories, postCategory)
	}

	return postCategories, nil
}
