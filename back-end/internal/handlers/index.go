package handlers

import (
	"database/sql"
	"net/http"

	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)

func Index(db *sql.DB, v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		categories, err := models.GetCategories(db)
		if err != nil {
			http.Error(w, "Erreur catégories", http.StatusInternalServerError)
			return
		}

		posts, err := models.GetHomePosts(db)
		if err != nil {
			http.Error(w, "Erreur posts", http.StatusInternalServerError)
			return
		}

		v.Render(w, "HomePage.html", map[string]any{
			"Categories": categories,
			"Posts":      posts,
		})
	}
}