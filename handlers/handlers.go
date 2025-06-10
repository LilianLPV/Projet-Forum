package handlers

import (
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("./template/*.html"))
)

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "accueil", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}
