package entity

// Category represents a category that can be assigned to posts
type Category struct {
	CategoryID int64  `json:"category_id"`
	Name       string `json:"name"`
}
