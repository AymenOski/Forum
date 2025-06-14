package usecase

import (
	"errors"
	"sync"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type PostRateLimiter struct {
	userLastPost map[uuid.UUID]time.Time
	mutex        sync.RWMutex
	limitTime    time.Duration
}

func NewPostRateLimiter() *PostRateLimiter {
	return &PostRateLimiter{
		userLastPost: make(map[uuid.UUID]time.Time),
		limitTime:    30 * time.Second,
	}
}

type PostService struct {
	postRepo          repository.PostRepository
	userRepo          repository.UserRepository
	categoryRepo      repository.CategoryRepository
	postAggregateRepo repository.PostAggregateRepository
	postReactionRepo  repository.PostReactionRepository
	sessionRepo       repository.UserSessionRepository
	rateLimiter       *PostRateLimiter
}

func NewPostService(postRepo *repository.PostRepository, userRepo *repository.UserRepository,
	categoryRepo *repository.CategoryRepository, postCategoryRepo *repository.PostAggregateRepository,
	postReactionRepo *repository.PostReactionRepository, sessionRepo *repository.UserSessionRepository, postRateLimit *PostRateLimiter,
) *PostService {
	return &PostService{
		postRepo:          *postRepo,
		userRepo:          *userRepo,
		categoryRepo:      *categoryRepo,
		postAggregateRepo: *postCategoryRepo,
		postReactionRepo:  *postReactionRepo,
		sessionRepo:       *sessionRepo,
		rateLimiter:       postRateLimit,
	}
}

func (r *PostRateLimiter) CanUserPost(userID uuid.UUID) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	lastPost, exists := r.userLastPost[userID]
	if !exists {
		return true
	}
	elapsed := time.Since(lastPost)

	if elapsed > time.Second*30 {
		return true
	}
	return false
}

func (ps *PostService) CreatePost(token string, content string, categoryIDs []*uuid.UUID) (*entity.Post, error) {
	// Get session and user from token
	session, err := ps.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	user, err := ps.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	canPost := ps.rateLimiter.CanUserPost(user.ID)
	if !canPost {
		return nil, errors.New("you can't create a post now, wait a bit")
	}

	// Validate content
	if content == "" {
		return nil, errors.New("post content cannot be empty")
	}
	if len(content) > 5000 {
		return nil, errors.New("post content too long (max: 5000 characters)")
	}

	// Validate categories
	if len(categoryIDs) <= 0 {
		return nil, errors.New("you have to select one category at least")
	}

	for _, categoryID := range categoryIDs {
		category, err := ps.categoryRepo.GetByID(*categoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("one or more categories does not exist")
		}
	}

	post := &entity.Post{
		UserID:    user.ID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	// Create and associate the categories to the post
	err = ps.postAggregateRepo.CreatePostWithCategories(post, categoryIDs)
	if err != nil {
		return nil, err
	}

	ps.rateLimiter.mutex.Lock()
	ps.rateLimiter.userLastPost[user.ID] = time.Now()
	ps.rateLimiter.mutex.Unlock()

	return post, nil
}

// ReactToPost - Like/dislike a post with toggle support.
// Same reaction twice = remove (toggle), different reaction = update.
// Returns the reaction entity on all operations (including delete for UI feedback).
// Parameters: postID, userID, reaction (true=like, false=dislike)
func (ps PostService) ReactToPost(postID uuid.UUID, token string, reaction bool) (*entity.PostReaction, error) {
	session, err := ps.sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}

	_, err = ps.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, err
	}

	_, err = ps.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	pr, err := ps.postReactionRepo.GetByUserAndPost(session.UserID, postID)
	if err == nil {
		if pr.Reaction == reaction {
			err := ps.postReactionRepo.Delete(pr.ID)
			if err != nil {
				return nil, errors.New("mistake in updating the post reaction")
			}
			return pr, nil
		} else if pr.Reaction != reaction {
			pr.Reaction = reaction
			pr.CreatedAt = time.Now()
			err := ps.postReactionRepo.Update(pr)
			if err != nil {
				return nil, errors.New("mistake in updating the post reaction")
			}
			return pr, nil
		}
	}
	PostReaction := &entity.PostReaction{
		UserID:    session.UserID,
		PostID:    postID,
		Reaction:  reaction,
		CreatedAt: time.Now(),
	}

	ps.postReactionRepo.Create(PostReaction)
	if err != nil {
		return nil, err
	}
	return PostReaction, nil
}

func (pc *PostService) GetPosts() ([]*entity.PostWithDetails, error) {
	posts, err := pc.postAggregateRepo.GetFeedForUser()
	if err != nil {
		return nil, err
	}
	return posts, nil
}
