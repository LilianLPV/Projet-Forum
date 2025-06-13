package handlers

import (
	"html/template"
	"net/http"
)

var (
	templates = template.Must(template.ParseGlob("./template/*.html"))
)

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home", nil)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "register", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}

func CreatepostHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create-post", nil)
}