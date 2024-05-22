package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/ykkssyaa/Posts_Service/internal/models"
)

type CommentsPostgres struct {
	db *sqlx.DB
}

func NewCommentsPostgres(db *sqlx.DB) *CommentsPostgres {
	return &CommentsPostgres{db: db}
}

func (c CommentsPostgres) CreateComment(comment models.Comment) (models.Comment, error) {

	tx, err := c.db.Begin()
	if err != nil {
		return models.Comment{}, err
	}

	query := `INSERT INTO comments (content, author, post, reply_to) 
				VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	row := tx.QueryRow(query, comment.Content, comment.Author, comment.Post, comment.ReplyTo)
	if err := row.Scan(&comment.ID, &comment.CreatedAt); err != nil {
		tx.Rollback()
		return models.Comment{}, err
	}

	return comment, tx.Commit()

}

func (c CommentsPostgres) GetCommentsByPost(postId int) ([]*models.Comment, error) {

	query := `SELECT FROM comments WHERE post = $1 AND reply_to IS NULL`

	var comments []*models.Comment

	if err := c.db.Select(&comments, query, postId); err != nil {
		return nil, err
	}

	return comments, nil
}
