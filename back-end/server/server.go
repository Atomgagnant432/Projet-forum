package server

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

func Start() {
	// Routes
	http.HandleFunc("/", indexHandler)

	// Static files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("front-end/static/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("front-end/static/img"))))

	fmt.Println("✅ Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
