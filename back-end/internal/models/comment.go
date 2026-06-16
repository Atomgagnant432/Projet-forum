package models

import (
	"database/sql"
)

type Comment struct {
	ID        int
	UserID    int
	Avatar    string
	Pseudo     string
	Content   string
	CreatedAt string
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]Comment, error) {
	rows, err := db.Query(`
		SELECT
			c.id,
			c.user_id,
			u.pseudo,
			c.content,
			c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment

		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.Pseudo,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		comment.Avatar = "👤"
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func InsertComment(db *sql.DB, postID int, userID int, content string) error {

	_, err := db.Exec(`
		INSERT INTO comments (
			post_id,
			user_id,
			content,
			created_at
		)
		VALUES (?, ?, ?, datetime('now'))
	`,
		postID,
		userID,
		content,
	)


	return err
}