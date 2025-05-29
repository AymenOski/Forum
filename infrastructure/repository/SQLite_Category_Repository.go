package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLiteCategoryRepository implements CategoryRepository interface
type SQLiteCategoryRepository struct {
	db *sql.DB
}

func NewSQLiteCategoryRepository(db *sql.DB) repository.CategoryRepository {
	return &SQLiteCategoryRepository{db: db}
}

func (r *SQLiteCategoryRepository) Create(category *entity.Category) error {
	category.ID = uuid.New()
	category.CreatedAt = time.Now()

	query := `INSERT INTO categories (id, name, description, created_at)
			  VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, category.ID.String(), category.Name, category.Description, category.CreatedAt)
	return err
}

func (r *SQLiteCategoryRepository) GetByID(categoryID uuid.UUID) (*entity.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories WHERE id = ?`

	row := r.db.QueryRow(query, categoryID.String())

	category := &entity.Category{}
	var idStr string

	err := row.Scan(&idStr, &category.Name, &category.Description, &category.CreatedAt)
	if err != nil {
		return nil, err
	}

	category.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *SQLiteCategoryRepository) GetByName(name string) (*entity.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories WHERE name = ?`

	row := r.db.QueryRow(query, name)

	category := &entity.Category{}
	var idStr string

	err := row.Scan(&idStr, &category.Name, &category.Description, &category.CreatedAt)
	if err != nil {
		return nil, err
	}

	category.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *SQLiteCategoryRepository) GetAll() ([]*entity.Category, error) {
	query := `SELECT id, name, description, created_at FROM categories ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category

	for rows.Next() {
		category := &entity.Category{}
		var idStr string

		err := rows.Scan(&idStr, &category.Name, &category.Description, &category.CreatedAt)
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

func (r *SQLiteCategoryRepository) Update(category *entity.Category) error {
	query := `UPDATE categories SET name = ?, description = ? WHERE id = ?`

	_, err := r.db.Exec(query, category.Name, category.Description, category.ID.String())
	return err
}

func (r *SQLiteCategoryRepository) Delete(categoryID uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = ?`

	_, err := r.db.Exec(query, categoryID.String())
	return err
}

func (r *SQLiteCategoryRepository) CheckNameExists(name string) (bool, error) {
	query := `SELECT COUNT(*) FROM categories WHERE name = ?`

	var count int
	err := r.db.QueryRow(query, name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *SQLiteCategoryRepository) GetWithPostCount() ([]*entity.Category, error) {
	query := `SELECT c.id, c.name, c.description, c.created_at, COUNT(pc.post_id) as post_count
			  FROM categories c 
			  LEFT JOIN post_categories pc ON c.id = pc.category_id 
			  GROUP BY c.id, c.name, c.description, c.created_at 
			  ORDER BY c.name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entity.Category

	for rows.Next() {
		category := &entity.Category{}
		var idStr string
		var postCount int // We're not storing this in the entity, but you could extend it

		err := rows.Scan(&idStr, &category.Name, &category.Description, &category.CreatedAt, &postCount)
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
