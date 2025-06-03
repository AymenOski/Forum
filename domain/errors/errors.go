package custom_errors

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrDatabaseError    = errors.New("database operation failed")
)

var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrPostNotFound    = errors.New("post not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidComment  = errors.New("invalid comment")
)

var (
	ErrReactionNotFound    = errors.New("reaction not found")
	ErrReactionExists      = errors.New("reaction already exists")
	ErrInvalidReactionType = errors.New("invalid reaction type")
)
