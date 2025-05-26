package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"

	"github.com/google/uuid"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(user *entity.User) error {
	// Generate UUID for new user
	user.UserID = uuid.New()

	query := `INSERT INTO users (user_id, name, email, password_hash, created_at) 
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, user.UserID, user.Name, user.Email, user.PasswordHash, time.Now())
	return err
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*entity.User, error) {
	query := `SELECT user_id, name, email, password_hash, session_token, session_expiry, created_at 
			  FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	user := &entity.User{}
	var userIDStr string
	var sessionToken sql.NullString
	var sessionExpiry sql.NullTime

	err := row.Scan(&userIDStr, &user.Name, &user.Email, &user.PasswordHash,
		&sessionToken, &sessionExpiry, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	user.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	if sessionToken.Valid {
		user.SessionToken = &sessionToken.String
	}
	if sessionExpiry.Valid {
		user.SessionExpiry = &sessionExpiry.Time
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetByID(userID *uuid.UUID) (*entity.User, error) {
	query := `SELECT user_id, name, email, password_hash, session_token, session_expiry, created_at 
			  FROM users WHERE user_id = ?`

	row := r.db.QueryRow(query, userID.String())

	user := &entity.User{}
	var userIDStr string
	var sessionToken sql.NullString
	var sessionExpiry sql.NullTime

	err := row.Scan(&userIDStr, &user.Name, &user.Email, &user.PasswordHash,
		&sessionToken, &sessionExpiry, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	user.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	if sessionToken.Valid {
		user.SessionToken = &sessionToken.String
	}
	if sessionExpiry.Valid {
		user.SessionExpiry = &sessionExpiry.Time
	}

	return user, nil
}

func (r *SQLiteUserRepository) GetByUsername(username string) (*entity.User, error) {
	// Assuming you have a username field, or use name field
	query := `SELECT user_id, name, email, password_hash, session_token, session_expiry, created_at 
			  FROM users WHERE name = ?`

	row := r.db.QueryRow(query, username)

	user := &entity.User{}
	var userIDStr string
	var sessionToken sql.NullString
	var sessionExpiry sql.NullTime

	err := row.Scan(&userIDStr, &user.Name, &user.Email, &user.PasswordHash,
		&sessionToken, &sessionExpiry, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	user.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	if sessionToken.Valid {
		user.SessionToken = &sessionToken.String
	}
	if sessionExpiry.Valid {
		user.SessionExpiry = &sessionExpiry.Time
	}

	return user, nil
}

func (r *SQLiteUserRepository) CreateSession(userID *uuid.UUID, sessionToken string, expiry time.Time) error {
	query := `UPDATE users SET session_token = ?, session_expiry = ? WHERE user_id = ?`
	_, err := r.db.Exec(query, sessionToken, expiry, userID.String())
	return err
}

func (r *SQLiteUserRepository) GetBySessionToken(sessionToken string) (*entity.User, error) {
	query := `SELECT user_id, name, email, password_hash, session_token, session_expiry, created_at 
			  FROM users WHERE session_token = ? AND session_expiry > ?`

	row := r.db.QueryRow(query, sessionToken, time.Now())

	user := &entity.User{}
	var userIDStr string
	var sessionTokenNull sql.NullString
	var sessionExpiry sql.NullTime

	err := row.Scan(&userIDStr, &user.Name, &user.Email, &user.PasswordHash,
		&sessionTokenNull, &sessionExpiry, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Parse UUID
	user.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	if sessionTokenNull.Valid {
		user.SessionToken = &sessionTokenNull.String
	}
	if sessionExpiry.Valid {
		user.SessionExpiry = &sessionExpiry.Time
	}

	return user, nil
}

func (r *SQLiteUserRepository) ClearSession(userID *uuid.UUID) error {
	query := `UPDATE users SET session_token = NULL, session_expiry = NULL WHERE user_id = ?`
	_, err := r.db.Exec(query, userID.String())
	return err
}

func (r *SQLiteUserRepository) CleanExpiredSessions() error {
	query := `UPDATE users SET session_token = NULL, session_expiry = NULL 
			  WHERE session_expiry < ?`
	_, err := r.db.Exec(query, time.Now())
	return err
}
