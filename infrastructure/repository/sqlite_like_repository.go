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
	tableName := "comment_reactions"
	if isPost {
		tableName = "post_reactions"
	}

	// Check if reaction already exists
	var existingReaction bool
	query := fmt.Sprintf(`SELECT EXISTS(
		SELECT 1 FROM %s 
		WHERE %s_id = ? AND user_id = ?
	)`, tableName, getParentColumn(isPost))
	
	err := r.db.QueryRow(query, parentID, userID).Scan(&existingReaction)
	if err != nil {
		return fmt.Errorf("failed to check existing reaction: %w", err)
	}

	if existingReaction {
		// Update existing reaction
		query = fmt.Sprintf(`UPDATE %s SET is_like = ? 
				WHERE %s_id = ? AND user_id = ?`, tableName, getParentColumn(isPost))
		_, err = r.db.Exec(query, isLike, parentID, userID)
	} else {
		// Insert new reaction
		query = fmt.Sprintf(`INSERT INTO %s (%s_id, user_id, is_like)
				VALUES (?, ?, ?)`, tableName, getParentColumn(isPost))
		_, err = r.db.Exec(query, parentID, userID, isLike)
	}

	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

func (r *SQLiteLikeRepository) RemoveParentReaction(parentID int64, userID uuid.UUID, isPost bool) error {
	tableName := "comment_reactions"
	if isPost {
		tableName = "post_reactions"
	}

	query := fmt.Sprintf(`DELETE FROM %s 
			WHERE %s_id = ? AND user_id = ?`, tableName, getParentColumn(isPost))
	result, err := r.db.Exec(query, parentID, userID)
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
	tableName := "comment_reactions"
	if isPost {
		tableName = "post_reactions"
	}

	query := fmt.Sprintf(`SELECT id, %s_id, user_id, is_like
			FROM %s
			WHERE %s_id = ? AND user_id = ?`, getParentColumn(isPost), tableName, getParentColumn(isPost))
	row := r.db.QueryRow(query, parentID, userID)

	reaction := &entity.LikeDislike{}
	var userIDStr string

	err := row.Scan(
		&reaction.ID,
		&reaction.ParentID,
		&userIDStr,
		&reaction.IsLike,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No reaction exists (not an error)
		}
		return nil, fmt.Errorf("failed to get reaction: %w", err)
	}

	reaction.UserID = userIDStr
	reaction.IsPost = isPost
	return reaction, nil
}

func (r *SQLiteLikeRepository) CountParentLikes(parentID int64, isPost bool) (int, error) {
	tableName := "comment_reactions"
	if isPost {
		tableName = "post_reactions"
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s
			WHERE %s_id = ? AND is_like = ?`, tableName, getParentColumn(isPost))
	var count int
	err := r.db.QueryRow(query, parentID, true).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count likes: %w", err)
	}
	return count, nil
}

func (r *SQLiteLikeRepository) CountParentDislikes(parentID int64, isPost bool) (int, error) {
	tableName := "comment_reactions"
	if isPost {
		tableName = "post_reactions"
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s
			WHERE %s_id = ? AND is_like = ?`, tableName, getParentColumn(isPost))
	var count int
	err := r.db.QueryRow(query, parentID, false).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count dislikes: %w", err)
	}
	return count, nil
}

// Helper function to get the correct column name based on reaction type
func getParentColumn(isPost bool) string {
	if isPost {
		return "post"
	}
	return "comment"
}