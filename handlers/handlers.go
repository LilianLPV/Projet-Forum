package handlers

import (
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("../templates/*.html"))
)

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "accueil", nil)
}
