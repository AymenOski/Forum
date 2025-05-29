package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type CommentReactionRepository interface {
	Create(reaction *entity.CommentReaction) error
	GetByID(reactionID uuid.UUID) (*entity.CommentReaction, error)
	GetByUserAndComment(userID, commentID uuid.UUID) (*entity.CommentReaction, error)
	GetByCommentID(commentID uuid.UUID) ([]*entity.CommentReaction, error)
	GetByUserID(userID uuid.UUID) ([]*entity.CommentReaction, error)
	Update(reaction *entity.CommentReaction) error
	Delete(reactionID uuid.UUID) error
	DeleteByUserAndComment(userID, commentID uuid.UUID) error
	GetLikeCountByCommentID(commentID uuid.UUID) (int, error)
	GetDislikeCountByCommentID(commentID uuid.UUID) (int, error)
	GetReactionCountsByCommentID(commentID uuid.UUID) (likes int, dislikes int, err error)
	HasUserReacted(userID, commentID uuid.UUID) (bool, *bool, error) // exists, reaction_value, error
}
