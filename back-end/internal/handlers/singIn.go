package handlers


import (
	"database/sql"
	"net/http"
	"strings"
	"html/template"
	"forum/back-end/internal/models"


	"golang.org/x/crypto/bcrypt"
)

func renderLoginPage(w http.ResponseWriter, errorMessage, email string) {
	t, err := template.ParseFiles("front-end/template/SignIn.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := AuthPageData{
		ErrorMessage: errorMessage,
		Email:        email,
	}

	t.Execute(w, data)
}


func HandlerConnexion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			renderLoginPage(w, "", "")
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("pwd")

		row := db.QueryRow(`SELECT id, password_hash, pseudo FROM users WHERE email = ?`, email,)

		var id int
		var hash string
		var pseudo string

		err := row.Scan(&id, &hash, &pseudo)

		if err == sql.ErrNoRows {
			renderLoginPage(w, "Identifiants incorrects", email)
			return
		}
		if err != nil {
			renderLoginPage(w, "Erreur base de données", email)
			return
		}

		res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if res != nil {
			renderLoginPage(w, "Identifiants incorrects", email)
			return
		}
		
		err = models.CreateSession(w, id)
		if err != nil {
			http.Error(w, "Erreur création session", http.StatusInternalServerError)
			return
		}
		

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}