package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLiteCommentReactionRepository implements CommentReactionRepository interface
type SQLiteCommentReactionRepository struct {
	db *sql.DB
}

func NewSQLiteCommentReactionRepository(db *sql.DB) repository.CommentReactionRepository {
	return &SQLiteCommentReactionRepository{db: db}
}

func (r *SQLiteCommentReactionRepository) Create(reaction *entity.CommentReaction) error {
	reaction.ID = uuid.New()
	reaction.CreatedAt = time.Now()

	query := `INSERT INTO comment_reaction (id, user_id, comment_id, reaction, created_at)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, reaction.ID.String(), reaction.UserID.String(),
		reaction.CommentID.String(), reaction.Reaction, reaction.CreatedAt)
	return err
}

func (r *SQLiteCommentReactionRepository) GetByID(reactionID uuid.UUID) (*entity.CommentReaction, error) {
	query := `SELECT id, user_id, comment_id, reaction, created_at FROM comment_reaction WHERE id = ?`

	row := r.db.QueryRow(query, reactionID.String())

	reaction := &entity.CommentReaction{}
	var idStr, userIDStr, commentIDStr string

	err := row.Scan(&idStr, &userIDStr, &commentIDStr, &reaction.Reaction, &reaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	reaction.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	reaction.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	reaction.CommentID, err = uuid.Parse(commentIDStr)
	if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *SQLiteCommentReactionRepository) GetByUserAndComment(userID, commentID uuid.UUID) (*entity.CommentReaction, error) {
	query := `SELECT id, user_id, comment_id, reaction, created_at 
			  FROM comment_reaction WHERE user_id = ? AND comment_id = ?`

	row := r.db.QueryRow(query, userID.String(), commentID.String())

	reaction := &entity.CommentReaction{}
	var idStr, userIDStr, commentIDStr string

	err := row.Scan(&idStr, &userIDStr, &commentIDStr, &reaction.Reaction, &reaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	reaction.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	reaction.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	reaction.CommentID, err = uuid.Parse(commentIDStr)
	if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *SQLiteCommentReactionRepository) GetByCommentID(commentID uuid.UUID) ([]*entity.CommentReaction, error) {
	query := `SELECT id, user_id, comment_id, reaction, created_at 
			  FROM comment_reaction WHERE comment_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, commentID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*entity.CommentReaction

	for rows.Next() {
		reaction := &entity.CommentReaction{}
		var idStr, userIDStr, commentIDStr string

		err := rows.Scan(&idStr, &userIDStr, &commentIDStr, &reaction.Reaction, &reaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		reaction.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		reaction.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		reaction.CommentID, err = uuid.Parse(commentIDStr)
		if err != nil {
			return nil, err
		}

		reactions = append(reactions, reaction)
	}

	return reactions, nil
}

func (r *SQLiteCommentReactionRepository) GetByUserID(userID uuid.UUID) ([]*entity.CommentReaction, error) {
	query := `SELECT id, user_id, comment_id, reaction, created_at 
			  FROM comment_reaction WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*entity.CommentReaction

	for rows.Next() {
		reaction := &entity.CommentReaction{}
		var idStr, userIDStr, commentIDStr string

		err := rows.Scan(&idStr, &userIDStr, &commentIDStr, &reaction.Reaction, &reaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		reaction.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		reaction.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}

		reaction.CommentID, err = uuid.Parse(commentIDStr)
		if err != nil {
			return nil, err
		}

		reactions = append(reactions, reaction)
	}

	return reactions, nil
}

func (r *SQLiteCommentReactionRepository) Update(reaction *entity.CommentReaction) error {
	query := `UPDATE comment_reaction SET reaction = ? WHERE id = ?`

	_, err := r.db.Exec(query, reaction.Reaction, reaction.ID.String())
	return err
}

func (r *SQLiteCommentReactionRepository) Delete(reactionID uuid.UUID) error {
	query := `DELETE FROM comment_reaction WHERE id = ?`

	_, err := r.db.Exec(query, reactionID.String())
	return err
}

func (r *SQLiteCommentReactionRepository) DeleteByUserAndComment(userID, commentID uuid.UUID) error {
	query := `DELETE FROM comment_reaction WHERE user_id = ? AND comment_id = ?`

	_, err := r.db.Exec(query, userID.String(), commentID.String())
	return err
}

func (r *SQLiteCommentReactionRepository) GetLikeCountByCommentID(commentID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comment_reaction WHERE comment_id = ? AND reaction = 1`

	var count int
	err := r.db.QueryRow(query, commentID.String()).Scan(&count)
	return count, err
}

func (r *SQLiteCommentReactionRepository) GetDislikeCountByCommentID(commentID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comment_reaction WHERE comment_id = ? AND reaction = 0`

	var count int
	err := r.db.QueryRow(query, commentID.String()).Scan(&count)
	return count, err
}

func (r *SQLiteCommentReactionRepository) GetReactionCountsByCommentID(commentID uuid.UUID) (likes int, dislikes int, err error) {
	query := `SELECT 
				SUM(CASE WHEN reaction = 1 THEN 1 ELSE 0 END) as likes,
				SUM(CASE WHEN reaction = 0 THEN 1 ELSE 0 END) as dislikes
			  FROM comment_reaction WHERE comment_id = ?`

	err = r.db.QueryRow(query, commentID.String()).Scan(&likes, &dislikes)
	return likes, dislikes, err
}

func (r *SQLiteCommentReactionRepository) HasUserReacted(userID, commentID uuid.UUID) (bool, *bool, error) {
	query := `SELECT reaction FROM comment_reaction WHERE user_id = ? AND comment_id = ?`

	var reaction bool
	err := r.db.QueryRow(query, userID.String(), commentID.String()).Scan(&reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil // User hasn't reacted
		}
		return false, nil, err
	}

	return true, &reaction, nil // User has reacted, return the reaction value
}
