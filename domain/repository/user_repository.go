package repository

import (
	"forum/domain/entity"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(userID *uuid.UUID) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)

	CreateSession(userID *uuid.UUID, sessionToken string, expiry time.Time) error
	GetBySessionToken(sessionToken string) (*entity.User, error)
}