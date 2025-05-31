package custom_errors

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category already exists")
	ErrDatabaseError    = errors.New("database operation failed")
)