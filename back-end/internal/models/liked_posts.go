package models

import (
	"database/sql"
)

func GetLikedPostsByUser(db *sql.DB, userID string) ([]Post, error) {
	rows, err := db.Query(`
        SELECT
            p.id,
            p.user_id,
            p.title,
            p.content,
            p.image_url,
            p.created_at,
            u.pseudo,
            u.avatar,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND type = 'like')    AS likes,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND type = 'dislike') AS dislikes,
            1 AS user_liked,
            0 AS user_disliked
        FROM posts p
        JOIN likes l   ON l.post_id = p.id
        JOIN users u   ON u.id = p.user_id
        WHERE l.user_id = ? AND l.type = 'like'
        ORDER BY l.created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Title, &p.Content,
			&p.ImageLink, &p.CreatedAt,
			&p.Pseudo, &p.Avatar,
			&p.LikeCount, &p.DislikeCount,
			&p.UserLiked, &p.UserDisliked,
		)
		if err != nil {
			return nil, err
		}
		p.Categories, _ = GetCategoriesByPost(db, p.ID)
		posts = append(posts, p)
	}
	return posts, nil
}
