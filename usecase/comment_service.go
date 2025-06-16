package usecase

import (
	"errors"
	"strings"
	"sync"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type CommentRateLimiter struct {
	userLastComment map[uuid.UUID]time.Time
	mutex           sync.RWMutex
	limitTime       time.Duration
}

func NewCommentRateLimiter() *CommentRateLimiter {
	return &CommentRateLimiter{
		userLastComment: make(map[uuid.UUID]time.Time),
		limitTime:       30 * time.Second,
	}
}

type CommentService struct {
	userRepo            repository.UserRepository
	commentRepo         repository.CommentRepository
	postRepo            repository.PostRepository
	commentReactionRepo repository.CommentReactionRepository
	sessionRepo         repository.UserSessionRepository
	rateLimiter         *CommentRateLimiter
}

func NewCommentService(userRepo repository.UserRepository, commentRepo repository.CommentRepository,
	postRepo repository.PostRepository, sessionRepo repository.UserSessionRepository,
	commentReactionRepo repository.CommentReactionRepository, commentRateLimit *CommentRateLimiter,
) *CommentService {
	return &CommentService{
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		postRepo:            postRepo,
		commentReactionRepo: commentReactionRepo,
		sessionRepo:         sessionRepo,
		rateLimiter:         commentRateLimit,
	}
}

func (cs *CommentRateLimiter) CanUserComment(userID uuid.UUID) bool {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	lastComment, exists := cs.userLastComment[userID]
	if !exists {
		return true
	}
	elapsed := time.Since(lastComment)

	if elapsed > time.Second*30 {
		return true
	}
	return false
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

	canComment := cs.rateLimiter.CanUserComment(user.ID)
	if !canComment {
		return nil, errors.New("you can't create a comment now, wait a bit")
	}

	content = strings.TrimSpace(content)
	if len(content) > 100 {
		return nil, errors.New("comment length excceds 250 characters")
	} else if content == "" {
		return nil, errors.New("comment should have at least 1 character")
	}

	_, err = cs.postRepo.GetByID(*postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	comment := &entity.Comment{
		Content: content,
		UserID:  user.ID,
		PostID:  *postID,
	}

	err = cs.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	cs.rateLimiter.mutex.Lock()
	cs.rateLimiter.userLastComment[user.ID] = time.Now()
	cs.rateLimiter.mutex.Unlock()

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
