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
//	GetAll() ([]*entity.User, error)
//	Update(user *entity.User) error
//	Delete(userID uuid.UUID) error
	CheckEmailExists(email string) (bool, error)
	CheckUserNameExists(userName string) (bool, error)
}
