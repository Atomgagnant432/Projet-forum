package models

import (
	"database/sql"

	"forum/back-end/pkg/utils"
)

type Post struct {
	ID           int
	UserID       int
	Avatar       string
	Pseudo       string
	Title        string
	Content      string
	ImageLink    string
	CreatedAt    string
	CommentCount int
	Comments     []Comment
	LikeCount    int
	DislikeCount int
	UserLiked    int
	UserDisliked int
	Categories   []Category
}

func GetHomePosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query(`
		SELECT
			p.id,
			p.user_id,
			u.pseudo,
			p.content,
			COALESCE(p.image_link, ''),
			p.created_at,
			COUNT(DISTINCT c.id),
			COUNT(DISTINCT pl.id),
			COUNT(DISTINCT pd.id)
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN posts_like pl ON pl.post_id = p.id
		LEFT JOIN posts_dislike pd ON pd.post_id = p.id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Pseudo,
			&post.Content,
			&post.ImageLink,
			&post.CreatedAt,
			&post.CommentCount,
			&post.LikeCount,
			&post.DislikeCount,
			&post.UserLiked,
			&post.UserDisliked,
		)
		if err != nil {
			return nil, err
		}

		post.Avatar = "👤"
		post.Title = utils.MakeTitle(post.Content)

		comments, err := GetCommentsByPostID(db, post.ID)
		if err != nil {
			return nil, err
		}

		post.Comments = comments
		posts = append(posts, post)
	}

	return posts, rows.Err()
}
