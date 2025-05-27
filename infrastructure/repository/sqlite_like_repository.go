package infra_repository

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/domain/entity"

	"github.com/google/uuid"
)

type SQLiteLikeRepository struct {
	db *sql.DB
}

func NewSQLiteLikeRepository(db *sql.DB) *SQLiteLikeRepository {
	return &SQLiteLikeRepository{db: db}
}

func (r *SQLiteLikeRepository) AddParentReaction(parentID int64, userID uuid.UUID, isLike bool, isPost bool) error {
	// First check if reaction already exists
	var existingReaction bool
	query := `SELECT EXISTS(
		SELECT 1 FROM reactions 
		WHERE parent_id = ? AND user_id = ? AND is_post = ?
	)`
	err := r.db.QueryRow(query, parentID, userID, isPost).Scan(&existingReaction)
	if err != nil {
		return fmt.Errorf("failed to check existing reaction: %w", err)
	}

	if existingReaction {
		// Update existing reaction
		query = `UPDATE reactions SET is_like = ? 
				WHERE parent_id = ? AND user_id = ? AND is_post = ?`
		_, err = r.db.Exec(query, isLike, parentID, userID, isPost)
	} else {
		// Insert new reaction
		query = `INSERT INTO reactions (parent_id, user_id, is_like, is_post)
				VALUES (?, ?, ?, ?)`
		_, err = r.db.Exec(query, parentID, userID, isLike, isPost)
	}

	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

func (r *SQLiteLikeRepository) RemoveParentReaction(parentID int64, userID uuid.UUID, isPost bool) error {
	query := `DELETE FROM reactions 
			WHERE parent_id = ? AND user_id = ? AND is_post = ?`
	result, err := r.db.Exec(query, parentID, userID, isPost)
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no reaction found to remove")
	}

	return nil
}

func (r *SQLiteLikeRepository) GetParentReaction(parentID int64, userID uuid.UUID, isPost bool) (*entity.LikeDislike, error) {
	query := `SELECT id, parent_id, user_id, is_like, is_post
			FROM reactions
			WHERE parent_id = ? AND user_id = ? AND is_post = ?`
	row := r.db.QueryRow(query, parentID, userID, isPost)

	reaction := &entity.LikeDislike{}
	var userIDStr string

	err := row.Scan(
		&reaction.ID,
		&reaction.ParentID,
		&userIDStr,
		&reaction.IsLike,
		&reaction.IsPost,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No reaction exists (not an error)
		}
		return nil, fmt.Errorf("failed to get reaction: %w", err)
	}

	reaction.UserID = userIDStr
	return reaction, nil
}

func (r *SQLiteLikeRepository) CountParentLikes(parentID int64, isPost bool) (int, error) {
	query := `SELECT COUNT(*) FROM reactions
			WHERE parent_id = ? AND is_like = ? AND is_post = ?`
	var count int
	err := r.db.QueryRow(query, parentID, true, isPost).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count likes: %w", err)
	}
	return count, nil
}

func (r *SQLiteLikeRepository) CountParentDislikes(parentID int64, isPost bool) (int, error) {
	query := `SELECT COUNT(*) FROM reactions
			WHERE parent_id = ? AND is_like = ? AND is_post = ?`
	var count int
	err := r.db.QueryRow(query, parentID, false, isPost).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count dislikes: %w", err)
	}
	return count, nil
}