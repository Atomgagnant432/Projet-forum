package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"forum/back-end/internal/models"
)

func CreateComment(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		currentUser, err := models.GetCurrentUser(db, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulaire invalide", http.StatusBadRequest)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "Post invalide", http.StatusBadRequest)
			return
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		err = models.InsertComment(db, postID, currentUser.ID, content)
		if err != nil {
			http.Error(w, "Erreur création commentaire", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}