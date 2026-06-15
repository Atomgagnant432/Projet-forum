package handlers

import (
	"net/http"

	"forum/back-end/internal/database"
	"forum/back-end/internal/render"
	"forum/back-end/internal/server"
)

func LikedPosts(v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := server.GetCurrentUser(r)
		if err != nil || user == nil {
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
			return
		}

		posts, err := database.GetLikedPostsByUser(user.ID)
		if err != nil {
			http.Error(w, "Erreur interne", http.StatusInternalServerError)
			return
		}

		v.Render(w, "LikedPosts.html", map[string]any{
			"User":  user,
			"Posts": posts,
		})
	}
}
