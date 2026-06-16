package handlers

import (
	"database/sql"
	"net/http"
	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)

type ProfilePageData struct {
	User    *models.CurrentUser
	Modal   string
	Error   string
	Success string
}

func ProfileHandler(db *sql.DB,v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Récupérer l'utilisateur connecté
		user, err := models.GetCurrentUser(db, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Lire ?modal=edit depuis l'URL
		modal := r.URL.Query().Get("modal")

		data := ProfilePageData{
			User:  user,
			Modal: modal,
		}

		 v.Render(w, "ProfilePage.html", data)
	}
}