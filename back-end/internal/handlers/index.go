package handlers

import (
	"net/http"

	"database/sql"

	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)

func Index(db *sql.DB, v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		categories, err := models.GetCategories(db)
		if err != nil {
			http.Error(w, "Erreur chargement catégories", http.StatusInternalServerError)
			return
		}

		currentUserID := 1

		posts, err := models.GetHomePosts(db, currentUserID)
		if err != nil {
			http.Error(w, "Erreur chargement posts", http.StatusInternalServerError)
			return
		}

		data := models.HomePageData{
			Categories: categories,
			Posts:      posts,
		}

		v.Render(w, "HomePage.html", data)
	}
}
