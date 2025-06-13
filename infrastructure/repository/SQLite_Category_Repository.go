package infra_repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"forum/domain/entity"
	custom_errors "forum/domain/errors"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type SQLiteCategoryRepository struct {
	db *sql.DB
}

func NewSQLiteCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &SQLiteCategoryRepository{db: db}
}

func (r *SQLiteCategoryRepository) Create(category *entity.Category) error {
	category.ID = uuid.New()
	category.CreatedAt = time.Now()

	query := `INSERT INTO categories (id, name, created_at)
			  VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, category.ID.String(), category.Name, category.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return custom_errors.ErrCategoryExists
		}
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	return nil
}

func (r *SQLiteCategoryRepository) GetByID(categoryID uuid.UUID) (*entity.Category, error) {
	query := `SELECT id, name, created_at FROM categories WHERE id = ?`

	row := r.db.QueryRow(query, categoryID.String())

	category := &entity.Category{}
	var idStr string

	err := row.Scan(&idStr, &category.Name, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	category.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	return category, nil
}

func (r *SQLiteCategoryRepository) GetByName(name string) (*entity.Category, error) {
	query := `SELECT id, name, created_at FROM categories WHERE name = ?`

	row := r.db.QueryRow(query, name)

	category := &entity.Category{}
	var idStr string

	err := row.Scan(&idStr, &category.Name, &category.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	category.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	return category, nil
}

func (r *SQLiteCategoryRepository) GetAll() ([]*entity.Category, error) {
	query := `SELECT id, name, created_at FROM categories ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	var categories []*entity.Category

	for rows.Next() {
		category := &entity.Category{}
		var idStr string

		err := rows.Scan(&idStr, &category.Name, &category.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		category.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (r *SQLiteCategoryRepository) Update(category *entity.Category) error {
	query := `UPDATE categories SET name = ? WHERE id = ?`

	result, err := r.db.Exec(query, category.Name, category.ID.String())
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return custom_errors.ErrCategoryExists
		}
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return custom_errors.ErrCategoryNotFound
	}

	return nil
}

func (r *SQLiteCategoryRepository) Delete(categoryID uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = ?`

	result, err := r.db.Exec(query, categoryID.String())
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return custom_errors.ErrCategoryNotFound
	}

	return nil
}

func (r *SQLiteCategoryRepository) CheckNameExists(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM categories WHERE name = ?`

	var count int
	err := r.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}

	return count > 0, nil
}

func (r *SQLiteCategoryRepository) GetWithPostCount() ([]*entity.Category, error) {
	query := `SELECT c.id, c.name, c.created_at, COUNT(pc.post_id) as post_count
			  FROM categories c 
			  LEFT JOIN post_categories pc ON c.id = pc.category_id 
			  GROUP BY c.id, c.name, c.created_at 
			  ORDER BY c.name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
	}
	defer rows.Close()

	var categories []*entity.Category

	for rows.Next() {
		category := &entity.Category{}
		var idStr string
		var postCount int

		err := rows.Scan(&idStr, &category.Name, &category.CreatedAt, &postCount)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		category.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", custom_errors.ErrDatabaseError, err)
		}

		categories = append(categories, category)
	}

	return categories, nil
}
