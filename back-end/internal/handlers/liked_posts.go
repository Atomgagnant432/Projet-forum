package handlers

import (
	"database/sql"
	"net/http"

	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)

func LikedPosts(db *sql.DB, v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := models.GetCurrentUser(db, r)

		posts, err := models.GetLikedPostsByUser(db, userID)
		if err != nil {
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}

		v.Render(w, "LikedPosts.html", map[string]any{
			"Posts": posts,
		})
	}
}
