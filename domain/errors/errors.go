package custom_errors

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotFound    = errors.New("post not found")
	ErrInvalidComment  = errors.New("invalid comment")
)

var (
	ErrReactionNotFound    = errors.New("reaction not found")
	ErrReactionExists      = errors.New("reaction already exists")
	ErrInvalidReactionType = errors.New("invalid reaction type")
)

// Login Form Errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid login credentials")
	ErrMissingFields      = errors.New("email and password are required")
	ErrDatabaseError      = errors.New("database operation failed during login")
)

// Register Form Errors
var (
	ErrEmailAlreadyExists        = errors.New("email already registered")
	ErrNameTaken                 = errors.New("username already taken")
	ErrWeakPassword              = errors.New("password is too weak")
	ErrInvalidEmail              = errors.New("invalid email format")
	ErrMissingRegistrationFields = errors.New("name, email, and password are required")
	ErrRegistrationFailed        = errors.New("registration process failed")
)
