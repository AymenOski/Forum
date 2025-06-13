package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) repository.UserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(user *entity.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	
	query := `INSERT INTO user (id, user_name, email, password_hash, created_at)
			  VALUES (?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, user.ID.String(), user.UserName, user.Email, user.PasswordHash, user.CreatedAt)
	return err
}

func (r *SQLiteUserRepository) GetByID(userID uuid.UUID) (*entity.User, error) {
	query := `SELECT id, user_name, email, password_hash, created_at FROM user WHERE id = ?`
	
	row := r.db.QueryRow(query, userID.String())
	
	user := &entity.User{}
	var idStr string
	
	err := row.Scan(&idStr, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `SELECT id, user_name, email, password_hash, created_at FROM user WHERE email = ?`
	
	row := r.db.QueryRow(query, email)
	
	user := &entity.User{}
	var idStr string
	
	err := row.Scan(&idStr, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *SQLiteUserRepository) GetByUserName(userName string) (*entity.User, error) {
	query := `SELECT id, user_name, email, password_hash, created_at FROM user WHERE user_name = ?`
	
	row := r.db.QueryRow(query, userName)
	
	user := &entity.User{}
	var idStr string
	
	err := row.Scan(&idStr, &user.UserName, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	user.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *SQLiteUserRepository) CheckEmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM user WHERE email = ?`
	
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *SQLiteUserRepository) CheckUserNameExists(userName string) (bool, error) {
	query := `SELECT COUNT(*) FROM user WHERE user_name = ?`
	
	var count int
	err := r.db.QueryRow(query, userName).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}