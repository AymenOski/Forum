package repository

import (
	"forum/domain/entity"
)

type UserAggregateRepository interface {
	CreateUserSession(user *entity.User) (*entity.UserSession, error)
	AuthenticateUser(email, password string) (*entity.User, *entity.UserSession, error)
}
