package repository

import (
	"forum/domain/entity"
)

type UserAggregateRepository interface {
	//	GetUserWithStats(userID uuid.UUID) (*entity.UserWithStats, error)
	// GetUserActivity(userID uuid.UUID, limit int) (posts []*entity.Post, comments []*entity.Comment, err error)
	CreateUserSession(user *entity.User) (*entity.UserSession, error)
	AuthenticateUser(email, password string) (*entity.User, *entity.UserSession, error)
}
