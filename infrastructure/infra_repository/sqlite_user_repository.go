package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type sqliteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) repository.UserRepository {
	return &sqliteUserRepo{
		db: db,
	}
}

func (u *sqliteUserRepo) Create(user *entity.User) error {
	
	return nil
}

func (u *sqliteUserRepo) GetByID(userID *uuid.UUID) (*entity.User, error) {
	return nil, nil
}

func (u *sqliteUserRepo) GetByUsername(username string) (*entity.User, error) {
	return nil, nil
}
func (u *sqliteUserRepo) GetByEmail(email string) (*entity.User, error) { return nil, nil }
func (u *sqliteUserRepo) CreateSession(userID *uuid.UUID, sessionToken string, expiry time.Time) error {
	return nil
}

func (u *sqliteUserRepo) GetBySessionToken(sessionToken string) (*entity.User, error) {
	return nil, nil
}
