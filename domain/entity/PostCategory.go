package entity

// PostCategory represents the many-to-many relationship between posts and categories
type PostCategory struct {
	PostID     string `json:"post_id"`
	CategoryID int64  `json:"category_id"`
}
