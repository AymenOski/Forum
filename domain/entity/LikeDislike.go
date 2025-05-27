package entity

// LikeDislike represents a user's like or dislike on a post
type LikeDislike struct {
	ID     int64  `json:"id"`
	ParentID int64  `json:"post_id"`
	UserID string `json:"user_id"`
	IsLike bool   `json:"is_like"`
	IsPost bool `json:"is_post"`
}
