package entity

type Post struct {
	UserID        string `json:"user_id"`
	PostID        int64  `json:"post_id"`
	Authorname    string `json:"username"`
	Content       string `json:"content"`
	LikesCount    int    `json:"likes_count"`
	DislikesCount int    `json:"dislikes_count"`
}
