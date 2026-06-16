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

	//loading the server
	mux := http.NewServeMux()

	//loading all static file needed
	fs := http.FileServer(http.Dir("front-end/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//routes
	mux.HandleFunc("/Index", handlers.HandlerRegister(db))
	mux.HandleFunc("/register", handlers.HandlerRegister(db))
	mux.HandleFunc("/login", handlers.HandlerConnexion(db))
	mux.HandleFunc("/liked-posts", handlers.LikedPosts(db, v))
	mux.HandleFunc("/", handlers.HandlerIndex(db, v))
	mux.HandleFunc("/post/like", handlers.HandlerLikePost(db))
	mux.HandleFunc("/post/dislike", handlers.HandlerDislikePost(db))
	mux.HandleFunc("/logout", handlers.HandlerLogout())
	mux.HandleFunc("/profile", handlers.ProfileHandler(db))

	//server started
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
