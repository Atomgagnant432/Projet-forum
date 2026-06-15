package models

import (
	"database/sql"
)

type User struct {
	ID        int
	Pseudo    string
	Email     string
	Avatar    string
	CreatedAt string
}


func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `
		SELECT
			id,
			pseudo,
			email,
			avatar,
			created_at
		FROM users
		WHERE id = ?
	`

	var user User

	err := db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Pseudo,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}


	return &user, nil
}