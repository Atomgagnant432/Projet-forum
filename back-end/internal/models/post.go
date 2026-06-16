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

func GetHomePosts(db *sql.DB, currentUserID int) ([]Post, error) {
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
			COUNT(DISTINCT pd.id),
			EXISTS(SELECT 1 FROM posts_like ul WHERE ul.post_id = p.id AND ul.user_id = ?),
			EXISTS(SELECT 1 FROM posts_dislike ud WHERE ud.post_id = p.id AND ud.user_id = ?)
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN posts_like pl ON pl.post_id = p.id
		LEFT JOIN posts_dislike pd ON pd.post_id = p.id
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, currentUserID, currentUserID)
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

func GetPostsByCategoryIDs(db *sql.DB, categoryIDs []int, currentUserID int) ([]Post, error) {
	if len(categoryIDs) == 0 {
		return GetHomePosts(db, currentUserID)
	}

	query := `
        SELECT DISTINCT p.id
        FROM posts p
        JOIN post_categories pc ON pc.post_id = p.id
        WHERE pc.category_id IN (
    `

	args := []any{}

	for i, id := range categoryIDs {
		if i > 0 {
			query += ","
		}

		query += "?"
		args = append(args, id)
	}

	query += ") ORDER BY p.created_at DESC"

	return getPostsFromIDQuery(db, query, currentUserID, args...)
}

func GetPostsByUserID(db *sql.DB, userID int, currentUserID int) ([]Post, error) {
	return getPostsFromIDQuery(db, `
        SELECT id
        FROM posts
        WHERE user_id = ?
        ORDER BY created_at DESC
    `, currentUserID, userID)
}

func GetLikedPostsByUserID(db *sql.DB, userID int, currentUserID int) ([]Post, error) {
	return getPostsFromIDQuery(db, `
        SELECT post_id
        FROM posts_like
        WHERE user_id = ?
        ORDER BY post_id DESC
    `, currentUserID, userID)
}

func getPostsFromIDQuery(db *sql.DB, idQuery string, currentUserID int, args ...any) ([]Post, error) {
	rows, err := db.Query(idQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var postID int

		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}

		post, err := GetPostByID(db, postID, currentUserID)
		if err != nil {
			return nil, err
		}

		posts = append(posts, *post)
	}

	return posts, rows.Err()
}

func GetPostByID(db *sql.DB, postID int, currentUserID int) (*Post, error) {
	var post Post

	err := db.QueryRow(`
        SELECT
            p.id,
            p.user_id,
            u.pseudo,
            p.title,
            p.content,
            COALESCE(p.image_link, ''),
            p.created_at,
            COUNT(DISTINCT c.id),
            COUNT(DISTINCT pl.id),
            COUNT(DISTINCT pd.id),
            EXISTS(SELECT 1 FROM posts_like ul WHERE ul.post_id = p.id AND ul.user_id = ?),
            EXISTS(SELECT 1 FROM posts_dislike ud WHERE ud.post_id = p.id AND ud.user_id = ?)
        FROM posts p
        JOIN users u ON u.id = p.user_id
        LEFT JOIN comments c ON c.post_id = p.id
        LEFT JOIN posts_like pl ON pl.post_id = p.id
        LEFT JOIN posts_dislike pd ON pd.post_id = p.id
        WHERE p.id = ?
        GROUP BY p.id
    `, currentUserID, currentUserID, postID).Scan(
		&post.ID,
		&post.UserID,
		&post.Pseudo,
		&post.Title,
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

	comments, err := GetCommentsByPostID(db, post.ID)
	if err != nil {
		return nil, err
	}

	post.Comments = comments

	return &post, nil
}

func CreatePost(db *sql.DB, userID int, title string, content string, imageLink string, categoryIDs []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var imageValue any = nil
	if imageLink != "" {
		imageValue = imageLink
	}

	result, err := tx.Exec(`INSERT INTO posts (user_id, image_link, title, content) VALUES (?, ?, ?, ?)`, userID, imageValue, title, content)
	if err != nil {
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, categoryID := range categoryIDs {
		_, err := stmt.Exec(postID, categoryID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
