package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type UserSessionRepository interface {
	Create(session *entity.UserSession) error
	GetByToken(token string) (*entity.UserSession, error)
	GetByUserID(userID uuid.UUID) (*entity.UserSession, error)
	Update(session *entity.UserSession) error
	Delete(sessionID uuid.UUID) error
	DeleteAllUserSessions(userID uuid.UUID) error
}
