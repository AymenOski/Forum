package usecase

import (
	"errors"
	"strings"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type CommentService struct {
	userRepo    repository.UserRepository
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	// next field is temperoraly until we have a proper middleware
	sessionRepo         repository.UserSessionRepository
	commentReactionRepo repository.CommentReactionRepository
}

func NewCommentService(userRepo repository.UserRepository, commentRepo repository.CommentRepository,
	postRepo repository.PostRepository, sessionRepo repository.UserSessionRepository, commentReactionRepo repository.CommentReactionRepository,
) *CommentService {
	return &CommentService{
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		postRepo:            postRepo,
		sessionRepo:         sessionRepo,
		commentReactionRepo: commentReactionRepo,
	}
}

func (cs *CommentService) CreateComment(postID *uuid.UUID, userID *uuid.UUID, content string) (*entity.Comment, error) {
	content = strings.TrimSpace(content)
	if len(content) > 249 {
		return nil, errors.New("comment length excceds 250 characters")
	} else if content == "" {
		return nil, errors.New("comment should have at least 1 character")
	}

	_, err := cs.postRepo.GetByID(*postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	_, err = cs.userRepo.GetByID(*userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	comment := &entity.Comment{
		Content: content,
		UserID:  *userID,
		PostID:  *postID,
	}

	err = cs.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// ReactToComment - Like/dislike a comment with toggle support.
// Same reaction twice = remove (toggle), different reaction = update.
// Returns nil when reaction is removed, reaction entity when created/updated.
func (cs *CommentService) ReactToComment(commentID *uuid.UUID, userID *uuid.UUID, reaction bool) (*entity.CommentReaction, error) {
	_, err := cs.userRepo.GetByID(*userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	_, err = cs.commentRepo.GetByID(*commentID)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	cr, err := cs.commentReactionRepo.GetByUserAndComment(*userID, *commentID)
	if err == nil {
		// user reacted, should update the reaction
		if cr.Reaction == reaction {
			cs.commentReactionRepo.Delete(cr.ID)
			return cr, nil
		} else if cr.Reaction != reaction {
			cr.Reaction = reaction
			cr.CreatedAt = time.Now()
			cs.commentReactionRepo.Update(cr)
			return cr, nil
		}
	}
	// no reaction of the user on the post, need to create a reaction
	commentReaction := &entity.CommentReaction{
		UserID:    *userID,
		CommentID: *commentID,
		Reaction:  reaction,
		CreatedAt: time.Now(),
	}
	err = cs.commentReactionRepo.Create(commentReaction)
	if err != nil {
		return nil, err
	}
	return commentReaction, nil
}

// temperoraly until we have a proper middleware
func (s *CommentService) GetUserFromSessionToken(token string) (*entity.User, error) {
	session, err := s.sessionRepo.GetByToken(token)
	if err != nil || session == nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
