package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type CommentReactionRepository interface {
	Create(reaction *entity.CommentReaction) error
	GetByUserAndComment(userID, commentID uuid.UUID) (*entity.CommentReaction, error)
	Update(reaction *entity.CommentReaction) error
	Delete(reactionID uuid.UUID) error
	 GetReactionCountsByCommentID(commentID uuid.UUID) (likes int, dislikes int, err error)
}
