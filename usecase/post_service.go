package usecase

import (
	"errors"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo          repository.PostRepository
	userRepo          repository.UserRepository
	categoryRepo      repository.CategoryRepository
	postAggregateRepo repository.PostAggregateRepository
	postReactionRepo  repository.PostReactionRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository, postCategoryRepo repository.PostAggregateRepository,
	postReactionRepo repository.PostReactionRepository) *PostService {
	return &PostService{
		postRepo:          postRepo,
		userRepo:          userRepo,
		categoryRepo:      categoryRepo,
		postAggregateRepo: postCategoryRepo,
		postReactionRepo:  postReactionRepo,
	}
}

func (ps *PostService) CreatePost(userID *uuid.UUID, content string, categoryIDs []*uuid.UUID) (*entity.Post, error) {
	user, err := ps.userRepo.GetByID(*userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if content == "" {
		return nil, errors.New("post content cannot be emtpy")
	}
	if len(categoryIDs) <= 0 {
		return nil, errors.New("you have to select one category at least")
	}
	if len(content) > 5000 {
		return nil, errors.New("post content too long (max: 5000 character)")
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
		UserID:    *userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	err = ps.postRepo.Create(post)
	if err != nil {
		return nil, err
	}
	// Associate the categories to the post
	err = ps.postAggregateRepo.CreatePostWithCategories(post, categoryIDs)
	if err != nil {
		err := ps.postRepo.Delete(post.ID)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return post, nil
}

// ReactToPost - Like/dislike a post with toggle support.
// Same reaction twice = remove (toggle), different reaction = update.
// Returns the reaction entity on all operations (including delete for UI feedback).
// Parameters: postID, userID, reaction (true=like, false=dislike)
func (ps *PostService) ReactToPost(postID *uuid.UUID, userID *uuid.UUID, reaction bool) (*entity.PostReaction, error) {
	_, err := ps.userRepo.GetByID(*userID)
	if err != nil {
		return nil, err
	}

	_, err = ps.postRepo.GetByID(*postID)
	if err != nil {
		return nil, err
	}

	pr, err := ps.postReactionRepo.GetByUserAndPost(*userID, *postID)
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
	commentReaction := &entity.PostReaction{
		UserID:    *userID,
		PostID:    *postID,
		Reaction:  reaction,
		CreatedAt: time.Now(),
	}

	ps.postReactionRepo.Create(commentReaction)
	if err != nil {
		return nil, err
	}
	return commentReaction, nil
}
