package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.UserSessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.UserSessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *AuthService) Register(name, email, password string) (*entity.User, error) {
	if !isValidEmail(email) {
		return nil, errors.New("this is the correct email address format : example@domain.com")
	}

	// Check if user already exists
	email = strings.ToLower(strings.TrimSpace(email))

	// Use the new CheckEmailExists method instead of GetByEmail
	exists, err := s.userRepo.CheckEmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	if len(name) < 3 || len(name) > 9 {
		return nil, errors.New("name should be between 3 and 9 characters")
	}

	name = strings.TrimSpace(name)
	if !isValidName(name) {
		return nil, errors.New("name should contain only letters and numbers")
	}

	// Check if username already exists
	userNameExists, err := s.userRepo.CheckUserNameExists(name)
	if err != nil {
		return nil, err
	}
	if userNameExists {
		return nil, errors.New("username already exists")
	}

	password = strings.TrimSpace(password)
	if len(password) < 6 || len(password) > 64 {
		return nil, errors.New("password should be between 6 and 64 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserName:     name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		// Note: ID and CreatedAt will be set in the repository Create method
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, *entity.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !isValidEmail(email) {
		return "", nil, errors.New("this is the correct email address format : example@domain.com")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("incorrect password")
	}

	// Clean up expired sessions first
	s.sessionRepo.DeleteExpiredSessions()

	// Check if user already has an active session
	existingSession, err := s.sessionRepo.GetByUserID(user.ID)
	if err == nil && existingSession != nil {
		// Check if the existing session is still valid
		if existingSession.ExpiresAt.After(time.Now()) {
			// Extend the existing session
			existingSession.ExpiresAt = time.Now().Add(24 * time.Hour)
			err = s.sessionRepo.Update(existingSession)
			if err != nil {
				return "", nil, err
			}
			return existingSession.SessionToken, user, nil
		} else {
			// Delete expired session
			s.sessionRepo.Delete(existingSession.ID)
		}
	}

	// Generate new session token
	token, err := s.generateSessionToken()
	if err != nil {
		return "", nil, err
	}

	expiry := time.Now().Add(24 * time.Hour)

	// Create new session
	session := &entity.UserSession{
		UserID:       user.ID,
		SessionToken: token,
		ExpiresAt:    expiry,
		// Note: ID and CreatedAt will be set in the repository Create method
	}

	err = s.sessionRepo.Create(session)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	return s.sessionRepo.DeleteAllUserSessions(userID)
}

func (s *AuthService) LogoutByToken(token string) error {
	return s.sessionRepo.DeleteByToken(token)
}

func (s *AuthService) ValidateSession(token string) (*entity.UserSession, error) {
	session, err := s.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	if session.ExpiresAt.Before(time.Now()) {
		// Delete expired session
		s.sessionRepo.Delete(session.ID)
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (s *AuthService) RefreshSession(token string) (string, error) {
	session, err := s.ValidateSession(token)
	if err != nil {
		return "", err
	}

	// Extend session expiry
	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	err = s.sessionRepo.Update(session)
	if err != nil {
		return "", err
	}

	return session.SessionToken, nil
}

func (s *AuthService) CleanupExpiredSessions() error {
	return s.sessionRepo.DeleteExpiredSessions()
}

func (s *AuthService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
