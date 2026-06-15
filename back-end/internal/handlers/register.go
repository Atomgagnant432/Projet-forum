package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"html/template"

	"golang.org/x/crypto/bcrypt"
)


type AuthPageData struct {
	ErrorMessage string
	Email        string
	Pseudo       string
}

func renderRegisterPage(
	w http.ResponseWriter, errorMessage string, email string, pseudo string,) {
	t, err := template.ParseFiles("front-end/template/SignUp.html")
	if err != nil {
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError,
		)
		return
	}

	data := AuthPageData{
		ErrorMessage: errorMessage,
		Email:        email,
		Pseudo:       pseudo,
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(
			w,
			"Erreur lors de l'affichage",
			http.StatusInternalServerError,
		)
	}
}

// Check the size of email and pseudo 
func validateRegister(email string,pseudo string,password string) string{
	msg := ""

	if !(len(email) >= 8 ){
		msg = msg + "Votre email doit contenir au moins 8 caractères "
	}

	if !(len(pseudo) >= 3){
		msg = msg + "\nVotre pseudo doit contenir au moins 3 caractères "
	}

	msg = msg + StrongPassword(password)

	return msg 
}


// Check if the password has at least a MAJ and a size of 8
func StrongPassword(password string) string {
	msg := ""
	VerMaj := "\nVotre mot de passe doit contenir au moins une majuscule."
	
	if len(password) < 8 {
		msg = msg + "\nVotre mot de passe doit contenir au moins 8 caractères"
	}

	for i := 0; i < len(password); i++ {
		if password[i] >= 'A' && password[i] <= 'Z' {
			VerMaj = ""
		}
	}

	msg = msg + VerMaj

	return msg
}


func HandlerRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			renderRegisterPage(w, "", "", "")
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		pseudo := strings.TrimSpace(r.FormValue("pseudo"))
		password := r.FormValue("pwd")
		passwordConfirm := r.FormValue("pwd-confirm")

		message := validateRegister(email, pseudo, password)
		if message != "" {
			renderRegisterPage(w, message, email, pseudo)
			return
		}

		if password != passwordConfirm {
			renderRegisterPage(
				w,
				"Les mots de passe ne correspondent pas",
				email,
				pseudo,
			)
			return
		}

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			http.Error(
				w,
				"Erreur lors du traitement du mot de passe",
				http.StatusInternalServerError,
			)
			return
		}

		_, err = db.ExecContext(
			r.Context(),
			`INSERT INTO users (email, pseudo, password_hash)
			 VALUES (?, ?, ?)`,
			email,
			pseudo,
			string(hash),
		)
		if err != nil {
			renderRegisterPage(
				w,
				"Email ou pseudo déjà utilisé",
				email,
				pseudo,
			)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}