package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLitePostAggregateRepository implements PostAggregateRepository interface
type SQLitePostAggregateRepository struct {
	db               *sql.DB
	postRepo         repository.PostRepository
	categoryRepo     repository.CategoryRepository
	postCategoryRepo repository.PostCategoryRepository
	userRepo         repository.UserRepository
	reactionRepo     repository.PostReactionRepository
}

func NewSQLitePostAggregateRepository(
	db *sql.DB,
	postRepo repository.PostRepository,
	categoryRepo repository.CategoryRepository,
	postCategoryRepo repository.PostCategoryRepository,
	userRepo repository.UserRepository,
	reactionRepo repository.PostReactionRepository,
) repository.PostAggregateRepository {
	return &SQLitePostAggregateRepository{
		db:               db,
		postRepo:         postRepo,
		categoryRepo:     categoryRepo,
		postCategoryRepo: postCategoryRepo,
		userRepo:         userRepo,
		reactionRepo:     reactionRepo,
	}
}

func (r *SQLitePostAggregateRepository) CreatePostWithCategories(post *entity.Post, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create post
	err = r.postRepo.Create(post)
	if err != nil {
		return err
	}

	// Associate categories
	for _, categoryID := range categoryIDs {
		postCategory := &entity.PostCategory{
			PostID:     post.ID,
			CategoryID: categoryID,
		}
		err = r.postCategoryRepo.Create(postCategory)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLitePostAggregateRepository) GetPostWithAllDetails(postID uuid.UUID) (*entity.PostWithDetails, error) {
	post, err := r.postRepo.GetByID(postID)
	if err != nil {
		return nil, err
	}

	author, err := r.userRepo.GetByID(post.UserID)
	if err != nil {
		return nil, err
	}

//	categories, err := r.postCategoryRepo.GetCategoriesByPostID(postID)
	if err != nil {
		return nil, err
	}

	likes, dislikes, err := r.reactionRepo.GetReactionCountsByPostID(postID)
	if err != nil {
		return nil, err
	}

	return &entity.PostWithDetails{
		Post:         *post,
		Author:       *author,
	//	Categories:   *categories,
		LikeCount:    likes,
		DislikeCount: dislikes,
	}, nil
}

func (r *SQLitePostAggregateRepository) GetFeedForUser(userID uuid.UUID, limit, offset int) ([]*entity.PostWithDetails, error) {
	// Simple implementation - get all recent posts
	posts, err := r.postRepo.GetWithPagination(limit, offset)
	if err != nil {
		return nil, err
	}

	var feedPosts []*entity.PostWithDetails
	for _, post := range posts {
		postWithDetails, err := r.GetPostWithAllDetails(post.ID)
		if err != nil {
			continue // Skip posts with errors
		}
		feedPosts = append(feedPosts, postWithDetails)
	}

	return feedPosts, nil
}

func (r *SQLitePostAggregateRepository) GetTrendingPosts(since time.Time, limit int) ([]*entity.PostWithDetails, error) {
	query := `SELECT p.id, p.title, p.content, p.user_id, p.created_at 
			  FROM posts p 
			  LEFT JOIN (
				  SELECT post_id, COUNT(*) as like_count 
				  FROM post_reaction 
				  WHERE reaction = 1 AND created_at > ? 
				  GROUP BY post_id
			  ) lr ON p.id = lr.post_id 
			  WHERE p.created_at > ?
			  ORDER BY COALESCE(lr.like_count, 0) DESC, p.created_at DESC 
			  LIMIT ?`

	rows, err := r.db.Query(query, since, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trendingPosts []*entity.PostWithDetails

	for rows.Next() {
		post := &entity.Post{}
		var idStr, userIDStr string

		err := rows.Scan(&idStr, &post.Title, &post.Content, &userIDStr, &post.CreatedAt)
		if err != nil {
			continue
		}

		post.ID, err = uuid.Parse(idStr)
		if err != nil {
			continue
		}

		post.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			continue
		}

		postWithDetails, err := r.GetPostWithAllDetails(post.ID)
		if err != nil {
			continue
		}

		trendingPosts = append(trendingPosts, postWithDetails)
	}

	return trendingPosts, nil
}

// SQLiteUserAggregateRepository implements UserAggregateRepository interface
type SQLiteUserAggregateRepository struct {
	db          *sql.DB
	userRepo    repository.UserRepository
	sessionRepo repository.UserSessionRepository
	postRepo    repository.PostRepository
	commentRepo repository.CommentRepository
}

func NewSQLiteUserAggregateRepository(
	db *sql.DB,
	userRepo repository.UserRepository,
	sessionRepo repository.UserSessionRepository,
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
) repository.UserAggregateRepository {
	return &SQLiteUserAggregateRepository{
		db:          db,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

// func (r *SQLiteUserAggregateRepository) GetUserWithStats(userID uuid.UUID) (*entity.UserWithStats, error) {
// 	user, err := r.userRepo.GetByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Get post count
// 	userPosts, err := r.postRepo.GetByUserID(userID)
// 	postCount := 0
// 	if err == nil {
// 		postCount = len(userPosts)
// 	}

// 	commentCount, err := r.commentRepo.GetCountByUserID(userID)
// 	if err != nil {
// 		commentCount = 0
// 	}

// 	return &entity.UserWithStats{
// 		User:         *user,
// 		PostCount:    postCount,
// 		CommentCount: commentCount,
// 	}, nil
// }

func (r *SQLiteUserAggregateRepository) GetUserActivity(userID uuid.UUID, limit int) (posts []*entity.Post, comments []*entity.Comment, err error) {
	posts, err = r.postRepo.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	// Limit posts if needed
	if len(posts) > limit {
		posts = posts[:limit]
	}

	comments, err = r.commentRepo.GetByUserID(userID)
	if err != nil {
		return posts, nil, err
	}

	// Limit comments if needed
	if len(comments) > limit {
		comments = comments[:limit]
	}

	return posts, comments, nil
}

func (r *SQLiteUserAggregateRepository) CreateUserWithSession(user *entity.User) (*entity.UserSession, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = r.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

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

func (r *SQLiteUserAggregateRepository) AuthenticateUser(email, password string) (*entity.User, *entity.UserSession, error) {
	user, err := r.userRepo.GetByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	// You would normally verify the password hash here
	// For now, assuming password verification is done elsewhere

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

// SQLiteRepositories implements the Repositories aggregator interface
type SQLiteRepositories struct {
	db                  *sql.DB
	userRepo            repository.UserRepository
	userSessionRepo     repository.UserSessionRepository
	postRepo            repository.PostRepository
	commentRepo         repository.CommentRepository
	categoryRepo        repository.CategoryRepository
	postCategoryRepo    repository.PostCategoryRepository
	postReactionRepo    repository.PostReactionRepository
	commentReactionRepo repository.CommentReactionRepository
	postAggregateRepo   repository.PostAggregateRepository
	userAggregateRepo   repository.UserAggregateRepository
}

func NewSQLiteRepositories(db *sql.DB) repository.Repositories {
	userRepo := NewSQLiteUserRepository(db)
	userSessionRepo := NewSQLiteUserSessionRepository(db)
	postRepo := NewSQLitePostRepository(db)
	commentRepo := NewSQLiteCommentRepository(db)
	categoryRepo := NewSQLiteCategoryRepository(db)
	//postCategoryRepo := NewSQLitePostCategoryRepository(db)
	postReactionRepo := NewSQLitePostReactionRepository(db)
	commentReactionRepo := NewSQLiteCommentReactionRepository(db)

	//postAggregateRepo := NewSQLitePostAggregateRepository(db, postRepo, categoryRepo, postCategoryRepo, userRepo, postReactionRepo)
	userAggregateRepo := NewSQLiteUserAggregateRepository(db, userRepo, userSessionRepo, postRepo, commentRepo)

	return &SQLiteRepositories{
		db:                  db,
		userRepo:            userRepo,
		userSessionRepo:     userSessionRepo,
		postRepo:            postRepo,
		commentRepo:         commentRepo,
		categoryRepo:        categoryRepo,
		//postCategoryRepo:    postCategoryRepo,
		postReactionRepo:    postReactionRepo,
		commentReactionRepo: commentReactionRepo,
		//postAggregateRepo:   postAggregateRepo,
		userAggregateRepo:   userAggregateRepo,
	}
}

// Factory methods to access individual repositories
func (r *SQLiteRepositories) User() repository.UserRepository               { return r.userRepo }
func (r *SQLiteRepositories) UserSession() repository.UserSessionRepository { return r.userSessionRepo }
func (r *SQLiteRepositories) Post() repository.PostRepository               { return r.postRepo }
func (r *SQLiteRepositories) Comment() repository.CommentRepository         { return r.commentRepo }
func (r *SQLiteRepositories) Category() repository.CategoryRepository       { return r.categoryRepo }
func (r *SQLiteRepositories) PostCategory() repository.PostCategoryRepository {
	return r.postCategoryRepo
}

func (r *SQLiteRepositories) PostReaction() repository.PostReactionRepository {
	return r.postReactionRepo
}

func (r *SQLiteRepositories) CommentReaction() repository.CommentReactionRepository {
	return r.commentReactionRepo
}

func (r *SQLiteRepositories) PostAggregate() repository.PostAggregateRepository {
	return r.postAggregateRepo
}

func (r *SQLiteRepositories) UserAggregate() repository.UserAggregateRepository {
	return r.userAggregateRepo
}
