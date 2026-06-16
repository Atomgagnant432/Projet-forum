package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)

func HandlerIndex(db *sql.DB, v *render.Render) http.HandlerFunc {
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

		currentUserID := 0

		currentUser, err := models.GetCurrentUser(db, r)
		if err == nil {
			currentUserID = currentUser.ID
		}

		filterType := r.URL.Query().Get("type")

		var posts []models.Post

		switch filterType {
		case "category":
			categoryIDs := []int{}

			for _, value := range r.URL.Query()["category_id"] {
				id, err := strconv.Atoi(value)
				if err == nil {
					categoryIDs = append(categoryIDs, id)
				}
			}

			posts, err = models.GetPostsByCategoryIDs(db, categoryIDs, currentUserID)

		case "created":
			if currentUserID == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			posts, err = models.GetPostsByUserID(db, currentUserID, currentUserID)

		case "liked":
			if currentUserID == 0 {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			posts, err = models.GetLikedPostsByUserID(db, currentUserID, currentUserID)

		default:
			posts, err = models.GetHomePosts(db, currentUserID)
		}

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
