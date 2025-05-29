package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLitePostReactionRepository implements PostReactionRepository interface
type SQLitePostReactionRepository struct {
	db *sql.DB
}

func NewSQLitePostReactionRepository(db *sql.DB) repository.PostReactionRepository {
	return &SQLitePostReactionRepository{db: db}
}

func (r *SQLitePostReactionRepository) Create(reaction *entity.PostReaction) error {
	reaction.ID = uuid.New()
	reaction.CreatedAt = time.Now()

	query := `INSERT INTO post_reaction (id, user_id, post_id, reaction, created_at)
			  VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, reaction.ID.String(), reaction.UserID.String(),
		reaction.PostID.String(), reaction.Reaction, reaction.CreatedAt)
	return err
}

func (r *SQLitePostReactionRepository) GetByID(reactionID uuid.UUID) (*entity.PostReaction, error) {
	query := `SELECT id, user_id, post_id, reaction, created_at FROM post_reaction WHERE id = ?`

	row := r.db.QueryRow(query, reactionID.String())

	reaction := &entity.PostReaction{}
	var idStr, userIDStr, postIDStr string

	err := row.Scan(&idStr, &userIDStr, &postIDStr, &reaction.Reaction, &reaction.CreatedAt)
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

	reaction.PostID, err = uuid.Parse(postIDStr)
	if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *SQLitePostReactionRepository) GetByUserAndPost(userID, postID uuid.UUID) (*entity.PostReaction, error) {
	query := `SELECT id, user_id, post_id, reaction, created_at 
			  FROM post_reaction WHERE user_id = ? AND post_id = ?`

	row := r.db.QueryRow(query, userID.String(), postID.String())

	reaction := &entity.PostReaction{}
	var idStr, userIDStr, postIDStr string

	err := row.Scan(&idStr, &userIDStr, &postIDStr, &reaction.Reaction, &reaction.CreatedAt)
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

	reaction.PostID, err = uuid.Parse(postIDStr)
	if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *SQLitePostReactionRepository) GetByPostID(postID uuid.UUID) ([]*entity.PostReaction, error) {
	query := `SELECT id, user_id, post_id, reaction, created_at 
			  FROM post_reaction WHERE post_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, postID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*entity.PostReaction

	for rows.Next() {
		reaction := &entity.PostReaction{}
		var idStr, userIDStr, postIDStr string

		err := rows.Scan(&idStr, &userIDStr, &postIDStr, &reaction.Reaction, &reaction.CreatedAt)
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

		reaction.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}

		reactions = append(reactions, reaction)
	}

	return reactions, nil
}

func (r *SQLitePostReactionRepository) GetByUserID(userID uuid.UUID) ([]*entity.PostReaction, error) {
	query := `SELECT id, user_id, post_id, reaction, created_at 
			  FROM post_reaction WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reactions []*entity.PostReaction

	for rows.Next() {
		reaction := &entity.PostReaction{}
		var idStr, userIDStr, postIDStr string

		err := rows.Scan(&idStr, &userIDStr, &postIDStr, &reaction.Reaction, &reaction.CreatedAt)
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

		reaction.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}

		reactions = append(reactions, reaction)
	}

	return reactions, nil
}

func (r *SQLitePostReactionRepository) Update(reaction *entity.PostReaction) error {
	query := `UPDATE post_reaction SET reaction = ? WHERE id = ?`

	_, err := r.db.Exec(query, reaction.Reaction, reaction.ID.String())
	return err
}

func (r *SQLitePostReactionRepository) Delete(reactionID uuid.UUID) error {
	query := `DELETE FROM post_reaction WHERE id = ?`

	_, err := r.db.Exec(query, reactionID.String())
	return err
}

func (r *SQLitePostReactionRepository) DeleteByUserAndPost(userID, postID uuid.UUID) error {
	query := `DELETE FROM post_reaction WHERE user_id = ? AND post_id = ?`

	_, err := r.db.Exec(query, userID.String(), postID.String())
	return err
}

func (r *SQLitePostReactionRepository) GetLikeCountByPostID(postID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM post_reaction WHERE post_id = ? AND reaction = 1`

	var count int
	err := r.db.QueryRow(query, postID.String()).Scan(&count)
	return count, err
}

func (r *SQLitePostReactionRepository) GetDislikeCountByPostID(postID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM post_reaction WHERE post_id = ? AND reaction = 0`

	var count int
	err := r.db.QueryRow(query, postID.String()).Scan(&count)
	return count, err
}

func (r *SQLitePostReactionRepository) GetReactionCountsByPostID(postID uuid.UUID) (likes int, dislikes int, err error) {
	query := `SELECT 
				SUM(CASE WHEN reaction = 1 THEN 1 ELSE 0 END) as likes,
				SUM(CASE WHEN reaction = 0 THEN 1 ELSE 0 END) as dislikes
			  FROM post_reaction WHERE post_id = ?`

	err = r.db.QueryRow(query, postID.String()).Scan(&likes, &dislikes)
	return likes, dislikes, err
}

func (r *SQLitePostReactionRepository) HasUserReacted(userID, postID uuid.UUID) (bool, *bool, error) {
	query := `SELECT reaction FROM post_reaction WHERE user_id = ? AND post_id = ?`

	var reaction bool
	err := r.db.QueryRow(query, userID.String(), postID.String()).Scan(&reaction)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil // User hasn't reacted
		}
		return false, nil, err
	}

	return true, &reaction, nil // User has reacted, return the reaction value
}
