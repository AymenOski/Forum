package usecase

import (
	"errors"
	"regexp"
)

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	res, _ := regexp.MatchString(emailRegex, email)
	return res
}

func isValidName(name string) bool {
	nameRegex := `^[a-zA-Z0-9]+$`
	res, _ := regexp.MatchString(nameRegex, name)
	return res
}

// 🧾 Validation: Generic Input Validation
var (
	ErrInvalidInputFormat   = errors.New("invalid input format")
	ErrFieldTooShort        = errors.New("input is too short")
	ErrFieldTooLong         = errors.New("input exceeds maximum allowed length")
	ErrMissingFieldsGeneric = errors.New("required fields are missing")
)

// 📩 Post Errors
var (
	ErrPostNotFound     = errors.New("post not found")
	ErrEmptyPostContent = errors.New("post content is required")
	ErrPostTooShort     = errors.New("post content is too short")
	ErrPostTooLong      = errors.New("post content exceeds maximum length")
	ErrInvalidPostID    = errors.New("invalid post ID format")
	ErrPostCreation     = errors.New("failed to create post")
)

// 💬 Comment Errors
var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrEmptyComment    = errors.New("comment content is required")
	ErrCommentTooShort = errors.New("comment is too short")
	ErrCommentTooLong  = errors.New("comment exceeds maximum allowed length")
	ErrInvalidComment  = errors.New("invalid comment")
	ErrCommentCreation = errors.New("failed to create comment")
)

// 🏷️ Category Errors
var (
	ErrCategoryNotFound    = errors.New("category not found")
	ErrEmptyCategoryName   = errors.New("category name is required")
	ErrInvalidCategoryName = errors.New("invalid category name format")
)

// 👍 Reaction Errors
var (
	ErrReactionNotFound    = errors.New("reaction not found")
	ErrReactionExists      = errors.New("reaction already exists")
	ErrInvalidReactionType = errors.New("invalid reaction type")
)

// 🔐 Auth/Login Errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid login credentials")
	ErrMissingFields      = errors.New("email and password are required")
	ErrDatabaseError      = errors.New("database operation failed during login")
	ErrUnauthorizedAccess = errors.New("unauthorized access to this resource")
	ErrSessionExpired     = errors.New("session expired, please log in again")
	ErrUserBanned         = errors.New("user account is banned")
)

// 📝 Register Errors
var (
	ErrEmailAlreadyExists        = errors.New("email already registered")
	ErrNameTaken                 = errors.New("username already taken")
	ErrWeakPassword              = errors.New("password is too weak")
	ErrInvalidEmail              = errors.New("invalid email format")
	ErrMissingRegistrationFields = errors.New("name, email, and password are required")
	ErrRegistrationFailed        = errors.New("registration process failed")
)
