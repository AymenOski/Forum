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
	userRepo            repository.UserRepository
	commentRepo         repository.CommentRepository
	postRepo            repository.PostRepository
	commentReactionRepo repository.CommentReactionRepository
	sessionRepo         repository.UserSessionRepository
}

func NewCommentService(userRepo *repository.UserRepository, commentRepo *repository.CommentRepository,
	postRepo *repository.PostRepository, commentReactionRepo *repository.CommentReactionRepository, sessionRepo *repository.UserSessionRepository,
) *CommentService {
	return &CommentService{
		userRepo:            *userRepo,
		commentRepo:         *commentRepo,
		postRepo:            *postRepo,
		commentReactionRepo: *commentReactionRepo,
		sessionRepo:         *sessionRepo,
	}
}

func (cs *CommentService) CreateComment(postID *uuid.UUID, token, content string) (*entity.Comment, error) {
	session, err := cs.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	user, err := cs.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	content = strings.TrimSpace(content)
	if len(content) > 250 {
		return nil, errors.New("comment length excceds 250 characters")
	} else if content == "" {
		return nil, errors.New("comment should have at least 1 character")
	}

	_, err = cs.postRepo.GetByID(*postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	comment := &entity.Comment{
		Content:   content,
		UserID:    user.ID,
		PostID:    *postID,
		CreatedAt: time.Now(),
	}
	// fmt.Printf("comment := &entity.Comment{ %v \n",comment)
	err = cs.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// ReactToComment - Like/dislike a comment with toggle support.
// Same reaction twice = remove (toggle), different reaction = update.
// Returns nil when reaction is removed, reaction entity when created/updated.
func (cs *CommentService) ReactToComment(commentID *uuid.UUID, token string, reaction bool) (*entity.CommentReaction, error) {
	session, err := cs.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	_, err = cs.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	_, err = cs.commentRepo.GetByID(*commentID)
	if err != nil {
		return nil, errors.New("comment not found")
	}

	cr, err := cs.commentReactionRepo.GetByUserAndComment(session.UserID, *commentID)
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
		UserID:    session.UserID,
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
