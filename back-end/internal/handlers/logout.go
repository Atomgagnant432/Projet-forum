package handlers

import (
	"net/http"

	"forum/back-end/internal/models"
)

func HandlerLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		models.DeleteSession(w, r)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}