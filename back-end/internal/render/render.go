package render

import (
	"html/template"
	"net/http"
)

type Render struct {
	tpl *template.Template
}

func New(pattern string) (*Render, error) {
	t, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, err
	}
	return &Render{tpl: t}, nil
}

func (v *Render) Render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := v.tpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}