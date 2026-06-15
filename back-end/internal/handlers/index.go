package handlers

import (
	"net/http"

	"forum/back-end/internal/models"
	"forum/back-end/internal/render"
)


func Index(r *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {


		posts, err := models.GetAllPosts()

		if err != nil {
			http.Error(w,"Database error",500)
			return
		}


		data := map[string]any{
			"Posts": posts,
		}


		r.Render(
			w,
			"HomePage.html",
			data,
		)
	}
}