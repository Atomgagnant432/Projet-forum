package handlers

import (
	"net/http"

	"forum/back-end/internal/render"
)

func Index(v *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			v.Render(w, "error.html", map[string]any{
				"Code":    404,
				"Message": "Page introuvable",
			})
			return
		}
		v.Render(w, "index.html", map[string]any{})
	}
}