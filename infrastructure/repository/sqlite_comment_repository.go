package infra_repository

import (
	"database/sql"
	"fmt"
	"time"

	"forum/domain/entity"

	"github.com/google/uuid"
)

type SQLiteCommentRepository struct {
	db *sql.DB
}

func NewSQLiteCommentRepository(db *sql.DB) *SQLiteCommentRepository {
	return &SQLiteCommentRepository{db: db}
}

func (r *SQLiteCommentRepository) Create(comment *entity.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content, created_at) 
			  VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, comment.PostID, comment.UserID, comment.Content, time.Now())
	if err != nil {
		return err
	}

	// Get the auto-incremented comment ID
	commentID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.CommentID = commentID

	return nil
}

func (r *SQLiteCommentRepository) GetByID(commentID int64) (*entity.Comment, error) {
	query := `SELECT c.comment_id, c.post_id, c.user_id, u.name as author_name, 
			 c.content, c.likes_count, c.dislikes_count, c.created_at
			 FROM comments c
			 JOIN users u ON c.user_id = u.user_id
			 WHERE c.comment_id = ?`

	row := r.db.QueryRow(query, commentID)

	comment := &entity.Comment{}
	var (
		postIDStr string
		userIDStr string
	)

	err := row.Scan(
		&comment.CommentID,
		&postIDStr,
		&userIDStr,
		&comment.Authorname,
		&comment.Content,
		&comment.LikesCount,
		&comment.DislikesCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, err
	}

	// Parse UUIDs
	comment.PostID = postIDStr
	comment.UserID = userIDStr

	return comment, nil
}

func (r *SQLiteCommentRepository) GetByPostID(postID uuid.UUID) ([]*entity.Comment, error) {
	query := `SELECT c.comment_id, c.post_id, c.user_id, u.name as author_name, 
			 c.content, c.likes_count, c.dislikes_count, c.created_at
			 FROM comments c
			 JOIN users u ON c.user_id = u.user_id
			 WHERE c.post_id = ?`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		comment := &entity.Comment{}
		var (
			postIDStr string
			userIDStr string
		)

		err := rows.Scan(
			&comment.CommentID,
			&postIDStr,
			&userIDStr,
			&comment.Authorname,
			&comment.Content,
			&comment.LikesCount,
			&comment.DislikesCount,
		)
		if err != nil {
			return nil, err
		}

		comment.PostID = postIDStr
		comment.UserID = userIDStr

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *SQLiteCommentRepository) GetByUserID(userID uuid.UUID) ([]*entity.Comment, error) {
	query := `SELECT c.comment_id, c.post_id, c.user_id, u.name as author_name, 
			 c.content, c.likes_count, c.dislikes_count, c.created_at
			 FROM comments c
			 JOIN users u ON c.user_id = u.user_id
			 WHERE c.user_id = ?`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		comment := &entity.Comment{}
		var (
			postIDStr string
			userIDStr string
		)

		err := rows.Scan(
			&comment.CommentID,
			&postIDStr,
			&userIDStr,
			&comment.Authorname,
			&comment.Content,
			&comment.LikesCount,
			&comment.DislikesCount,
		)
		if err != nil {
			return nil, err
		}

		comment.PostID = postIDStr
		comment.UserID = userIDStr

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *SQLiteCommentRepository) Update(comment *entity.Comment, userID uuid.UUID) error {
    // First, verify that the comment exists and belongs to the user
    var dbUserID string
    err := r.db.QueryRow("SELECT user_id FROM comments WHERE comment_id = ?", comment.CommentID).Scan(&dbUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("comment not found")
        }
        return err
    }

    // Check if the requesting user is the comment author
    if dbUserID != userID.String() {
        return fmt.Errorf("unauthorized: only the comment author can update the comment")
    }

    query := `UPDATE comments 
             SET content = ?, likes_count = ?, dislikes_count = ?
             WHERE comment_id = ?`

    _, err = r.db.Exec(query, comment.Content, comment.LikesCount, comment.DislikesCount, comment.CommentID)
    return err
}

func (r *SQLiteCommentRepository) Delete(commentID int64, userID uuid.UUID) error {
    // First, get the comment's author ID
    var dbUserID string
    err := r.db.QueryRow("SELECT user_id FROM comments WHERE comment_id = ?", commentID).Scan(&dbUserID)
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("comment not found")
        }
        return err
    }

    // Check if the requesting user is the author OR an admin
    if dbUserID != userID.String()  {
        return fmt.Errorf("unauthorized: only the comment author or an admin can delete the comment")
    }

    query := `DELETE FROM comments WHERE comment_id = ?`
    _, err = r.db.Exec(query, commentID)
    return err
}