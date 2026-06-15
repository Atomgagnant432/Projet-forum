package server

import (
	"fmt"
	"log"
	"net/http"

	"forum/back-end/database"
	"forum/back-end/internal/handlers"
	"forum/back-end/internal/render"
)

func Start() {
	db, err := database.OpenDatabase()
	if err != nil {
		log.Fatal("Erreur ouverture BDD :", err)
	}
	defer db.Close()

	//future render
	v, err := render.New("front-end/template/*.html")
	if err != nil {
		log.Fatal(err)
	}
	print(v)

	//Starting Server

	mux := http.NewServeMux()
	
	//Getting all filed
	fs := http.FileServer(http.Dir("front-end/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//routes
	// mux.HandleFunc("/", handlers.index(v))
	mux.HandleFunc("/register", handlers.HandlerRegister(db))
	mux.HandleFunc("/login", handlers.HandlerConnexion(db))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		v.Render(w, "HomePage.html", nil)

	})

	//server started
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}