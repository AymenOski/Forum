package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entity.User) error
	GetByID(userID uuid.UUID) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUserName(userName string) (*entity.User, error)
}
