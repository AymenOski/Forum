package entity

type Comment struct {
	CommentID     int64  `json:"comment_id"`
	PostID        string `json:"post_id"`
	UserID        string `json:"user_id"`
	Authorname    string `json:"username"`
	Content       string `json:"content"`
	LikesCount    int    `json:"likes_count"`
	DislikesCount int    `json:"dislikes_count"`
}
