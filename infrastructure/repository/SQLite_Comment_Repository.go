package infra_repository

import (
	"database/sql"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

// SQLiteCommentRepository implements CommentRepository interface
type SQLiteCommentRepository struct {
	db *sql.DB
}

func NewSQLiteCommentRepository(db *sql.DB) repository.CommentRepository {
	return &SQLiteCommentRepository{db: db}
}

func (r *SQLiteCommentRepository) Create(comment *entity.Comment) error {
	comment.ID = uuid.New()
	comment.CreatedAt = time.Now()
	
	query := `INSERT INTO comments (id, content, user_id, post_id, createdat)
			  VALUES (?, ?, ?, ?, ?)`
	
	_, err := r.db.Exec(query, comment.ID.String(), comment.Content, 
					   comment.UserID.String(), comment.PostID.String(), comment.CreatedAt)
	return err
}

func (r *SQLiteCommentRepository) GetByID(commentID uuid.UUID) (*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat FROM comments WHERE id = ?`
	
	row := r.db.QueryRow(query, commentID.String())
	
	comment := &entity.Comment{}
	var idStr, userIDStr, postIDStr string
	
	err := row.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	
	comment.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	
	comment.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}
	
	comment.PostID, err = uuid.Parse(postIDStr)
	if err != nil {
		return nil, err
	}
	
	return comment, nil
}

func (r *SQLiteCommentRepository) GetByPostID(postID uuid.UUID) ([]*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat 
			  FROM comments WHERE post_id = ? ORDER BY createdat ASC`
	
	rows, err := r.db.Query(query, postID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []*entity.Comment
	
	for rows.Next() {
		comment := &entity.Comment{}
		var idStr, userIDStr, postIDStr string
		
		err := rows.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		comment.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		
		comment.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		
		comment.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}
		
		comments = append(comments, comment)
	}
	
	return comments, nil
}

func (r *SQLiteCommentRepository) GetByUserID(userID uuid.UUID) ([]*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat 
			  FROM comments WHERE user_id = ? ORDER BY createdat DESC`
	
	rows, err := r.db.Query(query, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []*entity.Comment
	
	for rows.Next() {
		comment := &entity.Comment{}
		var idStr, userIDStr, postIDStr string
		
		err := rows.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		comment.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		
		comment.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		
		comment.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}
		
		comments = append(comments, comment)
	}
	
	return comments, nil
}

func (r *SQLiteCommentRepository) GetByPostIDWithPagination(postID uuid.UUID, limit, offset int) ([]*entity.Comment, error) {
	query := `SELECT id, content, user_id, post_id, createdat 
			  FROM comments WHERE post_id = ? ORDER BY createdat ASC LIMIT ? OFFSET ?`
	
	rows, err := r.db.Query(query, postID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var comments []*entity.Comment
	
	for rows.Next() {
		comment := &entity.Comment{}
		var idStr, userIDStr, postIDStr string
		
		err := rows.Scan(&idStr, &comment.Content, &userIDStr, &postIDStr, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		comment.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		
		comment.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		
		comment.PostID, err = uuid.Parse(postIDStr)
		if err != nil {
			return nil, err
		}
		
		comments = append(comments, comment)
	}
	
	return comments, nil
}

func (r *SQLiteCommentRepository) Update(comment *entity.Comment) error {
	query := `UPDATE comments SET content = ? WHERE id = ?`
	
	_, err := r.db.Exec(query, comment.Content, comment.ID.String())
	return err
}

func (r *SQLiteCommentRepository) Delete(commentID uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = ?`
	
	_, err := r.db.Exec(query, commentID.String())
	return err
}

func (r *SQLiteCommentRepository) GetCountByPostID(postID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = ?`
	
	var count int
	err := r.db.QueryRow(query, postID.String()).Scan(&count)
	return count, err
}

func (r *SQLiteCommentRepository) GetCountByUserID(userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE user_id = ?`
	
	var count int
	err := r.db.QueryRow(query, userID.String()).Scan(&count)
	return count, err
}

func (r *SQLiteCommentRepository) GetWithDetails(commentID uuid.UUID) (*entity.CommentWithDetails, error) {
	comment, err := r.GetByID(commentID)
	if err != nil {
		return nil, err
	}
	
	return &entity.CommentWithDetails{
		Comment: *comment,
	}, nil
}

func (r *SQLiteCommentRepository) GetByPostIDWithDetails(postID uuid.UUID) ([]*entity.CommentWithDetails, error) {
	comments, err := r.GetByPostID(postID)
	if err != nil {
		return nil, err
	}
	
	var commentsWithDetails []*entity.CommentWithDetails
	for _, comment := range comments {
		commentsWithDetails = append(commentsWithDetails, &entity.CommentWithDetails{
			Comment: *comment,
		})
	}
	
	return commentsWithDetails, nil
}