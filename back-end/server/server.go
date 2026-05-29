package server

import (
	"fmt"
	"log"
	"net/http"
)

func Start() {
	//future render

	//Starting Server

	mux := http.NewServeMux()
	
	//Getting all filed
	fs := http.FileServer(http.Dir("front-end/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//routes
	mux.HandleFunc("/", handlers.Home(v))

	//server started
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}