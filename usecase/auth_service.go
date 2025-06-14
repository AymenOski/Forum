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

func (s *AuthService) Signup(name, email, password string) (*entity.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if !isValidEmail(email) {
		return nil, errors.New("invalid email format. Make sure it follows the pattern: name@domain.com")
	}

	_, err := s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	if len(name) < 3 || len(name) > 9 {
		return nil, errors.New("name should be between 3 and 9 characters")
	}
	name = strings.TrimSpace(name)
	if !isValidName(name) {
		return nil, errors.New("name should contain only letters and numbers")
	}

	password = strings.TrimSpace(password)
	if len(password) < 6 || len(password) > 64 {
		return nil, errors.New("password should be between 4 and 64 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UserName:     name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
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
		return "", nil, errors.New("invalid email format. Make sure it follows the pattern: name@domain.com")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, errors.New("incorrect password")
	}

	s.sessionRepo.DeleteAllUserSessions(user.ID)

	token, err := s.generateSessionToken()
	if err != nil {
		return "", nil, err
	}

	session := &entity.UserSession{
		UserID:       user.ID,
		SessionToken: token,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
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

func (s *AuthService) ValidateSession(token string) (*entity.UserSession, error) {
	session, err := s.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}
	return session, nil
}

func (s *AuthService) GetUserFromSessionToken(token string) (*entity.User, error) {
	session, err := s.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
