package entity

type CommentWithDetails struct {
	Comment
	Author       User `json:"author"`
	LikeCount    int  `json:"like_count"`
	DislikeCount int  `json:"dislike_count"`
}
