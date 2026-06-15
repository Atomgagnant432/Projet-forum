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
	//loading database with error catch
	db, err := database.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//render all the template
	v, err := render.New("front-end/template/*.html")
	if err != nil {
		log.Fatal(err)
	}

	//loading the server
	mux := http.NewServeMux()

	//loading all static file needed
	fs := http.FileServer(http.Dir("front-end/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	mux.HandleFunc("/", handlers.Index(db, v))

	//starting the server
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}