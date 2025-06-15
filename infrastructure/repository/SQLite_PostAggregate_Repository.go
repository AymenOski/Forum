package infra_repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLitePostAggregateRepository implements PostAggregateRepository interface
type SQLitePostAggregateRepository struct {
	db               *sql.DB
	postRepo         repository.PostRepository
	postCategoryRepo repository.PostCategoryRepository
	userRepo         repository.UserRepository
	reactionRepo     repository.PostReactionRepository
	commentRepo      repository.CommentRepository
}

func NewSQLitePostAggregateRepository(
	db *sql.DB,
	postRepo *repository.PostRepository,
	postCategoryRepo *repository.PostCategoryRepository,
	userRepo *repository.UserRepository,
	reactionRepo *repository.PostReactionRepository,
	commentRepo *repository.CommentRepository,
) repository.PostAggregateRepository {
	return &SQLitePostAggregateRepository{
		db:               db,
		postRepo:         *postRepo,
		postCategoryRepo: *postCategoryRepo,
		userRepo:         *userRepo,
		reactionRepo:     *reactionRepo,
		commentRepo:      *commentRepo,
	}
}

// CreatePostWithCategories creates a post and associates it with categories
func (r *SQLitePostAggregateRepository) CreatePostWithCategories(post *entity.Post, categoryIDs []*uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = r.postRepo.Create(post)
	if err != nil {
		return err
	}

	// Associate categories
	for _, categoryID := range categoryIDs {
		postCategory := &entity.PostCategory{
			PostID:     post.ID,
			CategoryID: *categoryID,
		}
		err = r.postCategoryRepo.Create(postCategory)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLitePostAggregateRepository) GetFeedForUser() ([]*entity.PostWithDetails, error) {
	posts, err := r.postRepo.GetAll()
	if err != nil {
		return nil, err
	}
	postWithDetails := make([]*entity.PostWithDetails, 0, len(posts))
	for _, post := range posts {
		p, err := r.GetPostWithAllDetails(post.ID)
		if err != nil {
			return nil, err
		}
		postWithDetails = append(postWithDetails, p)
	}
	return postWithDetails, nil
}

// GetPostWithAllDetails retrieves a post with author, categories, and reaction counts , and comments!!
func (r *SQLitePostAggregateRepository) GetPostWithAllDetails(postID uuid.UUID) (*entity.PostWithDetails, error) {
	post, err := r.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	author, err := r.userRepo.GetByID(post.UserID)
	if err != nil {
		return nil, err
	}

	categories, err := r.postCategoryRepo.GetCategoriesByPostID(postID)
	if err != nil {
		return nil, err
	}

	likes, dislikes, err := r.reactionRepo.GetReactionCountsByPostID(postID)
	if err != nil {
		fmt.Println("reaction count")
		return nil, err
	}

	comments, err := r.commentRepo.GetByPostIDWithDetails(postID)
	if err != nil {
		fmt.Println("comments error")
		return nil, err
	}

	return &entity.PostWithDetails{
		Post:         *post,
		Author:       *author,
		Comments:     comments,
		Categories:   categories,
		LikeCount:    likes,
		DislikeCount: dislikes,
	}, nil
}

// SQLiteUserAggregateRepository implements UserAggregateRepository interface
type SQLiteUserAggregateRepository struct {
	db          *sql.DB
	userRepo    repository.UserRepository
	sessionRepo repository.UserSessionRepository
}

func NewSQLiteUserAggregateRepository(
	db *sql.DB,
	userRepo repository.UserRepository,
	sessionRepo repository.UserSessionRepository,
) repository.UserAggregateRepository {
	return &SQLiteUserAggregateRepository{
		db:          db,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// CreateUserSession creates a new session for a user
func (r *SQLiteUserAggregateRepository) CreateUserSession(user *entity.User) (*entity.UserSession, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	session := &entity.UserSession{
		UserID:       user.ID,
		SessionToken: uuid.New().String(),            // Generate a session token
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 24 hour session
	}

	err = r.sessionRepo.Create(session)
	if err != nil {
		return nil, err
	}

	return session, tx.Commit()
}

// AuthenticateUser verifies user credentials and creates a session
func (r *SQLiteUserAggregateRepository) AuthenticateUser(email, password string) (*entity.User, *entity.UserSession, error) {
	user, err := r.userRepo.GetByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	// Password verification should be done in the use case layer
	// This just creates the session after authentication

	// Create new session
	session := &entity.UserSession{
		UserID:       user.ID,
		SessionToken: uuid.New().String(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	err = r.sessionRepo.Create(session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (r *SQLitePostAggregateRepository) GetPostsWithDetailsByUser(userID uuid.UUID) ([]*entity.PostWithDetails, error) {
	posts, err := r.postRepo.GetbyuserId(userID)
	if err != nil {
		return nil, err
	}

	var postsWithDetails []*entity.PostWithDetails
	for _, post := range posts {
		pwd, err := r.GetPostWithAllDetails(post.ID)
		if err != nil {
			return nil, err
		}
		postsWithDetails = append(postsWithDetails, pwd)
	}

	return postsWithDetails, nil
}
