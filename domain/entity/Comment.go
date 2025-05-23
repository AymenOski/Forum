package entity

type Comment struct {
	CommentID int64  `json:"comment_id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
}
