package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"forum/back-end/internal/models"
)

func HandlerLikePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "post_id invalide", http.StatusBadRequest)
			return
		}

		currentUserID := 1

		if err := models.ToggleLike(db, currentUserID, postID); err != nil {
			http.Error(w, "Erreur lors du like", http.StatusInternalServerError)
			return
		}

		redirectBack(w, r)
	}
}

func HandlerDislikePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, "post_id invalide", http.StatusBadRequest)
			return
		}

		currentUserID := 1

		if err := models.ToggleDislike(db, currentUserID, postID); err != nil {
			http.Error(w, "Erreur lors du dislike", http.StatusInternalServerError)
			return
		}

		redirectBack(w, r)
	}
}

func redirectBack(w http.ResponseWriter, r *http.Request) {
	ref := r.Referer()
	if ref == "" {
		ref = "/"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}
