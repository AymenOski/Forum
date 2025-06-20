package entity

type PostWithDetails struct {
	Post
	Author       User        `json:"author"`
	Categories   []*Category `json:"categories,omitempty"`
	Comments     []CommentWithDetails   `json:"comments,omitempty"`
	LikeCount    int         `json:"like_count"`
	DislikeCount int         `json:"dislike_count"`
}

