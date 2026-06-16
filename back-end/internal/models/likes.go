package models

import "database/sql"

func HasUserLikedPost(db *sql.DB, userID, postID int) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM posts_like WHERE user_id = ? AND post_id = ?
		)
	`, userID, postID).Scan(&exists)
	return exists, err
}

func HasUserDislikedPost(db *sql.DB, userID, postID int) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM posts_dislike WHERE user_id = ? AND post_id = ?
		)
	`, userID, postID).Scan(&exists)
	return exists, err
}

func ToggleLike(db *sql.DB, userID, postID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	alreadyLiked, err := HasUserLikedPost(db, userID, postID)
	if err != nil {
		return err
	}

	if alreadyLiked {
		if _, err := tx.Exec(`DELETE FROM posts_like WHERE user_id = ? AND post_id = ?`, userID, postID); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(`DELETE FROM posts_dislike WHERE user_id = ? AND post_id = ?`, userID, postID); err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO posts_like (user_id, post_id) VALUES (?, ?)`, userID, postID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func ToggleDislike(db *sql.DB, userID, postID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	alreadyDisliked, err := HasUserDislikedPost(db, userID, postID)
	if err != nil {
		return err
	}

	if alreadyDisliked {
		if _, err := tx.Exec(`DELETE FROM posts_dislike WHERE user_id = ? AND post_id = ?`, userID, postID); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(`DELETE FROM posts_like WHERE user_id = ? AND post_id = ?`, userID, postID); err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO posts_dislike (user_id, post_id) VALUES (?, ?)`, userID, postID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
