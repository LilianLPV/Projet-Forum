package controllers

import (
	"fmt"
	"forum/hash"
	"forum/models"
	"forum/services"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type AccountControllers struct {
	service  *services.AccountService
	template *template.Template
}

func InitAccountController(service *services.AccountService, template *template.Template) *AccountControllers {
	return &AccountControllers{service: service, template: template}
}

func (c *AccountControllers) AccountRouter(r *mux.Router) {

	r.HandleFunc("/login", c.LoginForm).Methods("GET")
	r.HandleFunc("/login/treatment", c.Login).Methods("POST")
	r.HandleFunc("/signup", c.CreateForm).Methods("GET")
	r.HandleFunc("/signup/treatment", c.Create).Methods("POST")

	r.HandleFunc("/", c.AuthMiddleware(c.Home)).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(c.NotFound)
	r.HandleFunc("/unauthorized", c.Unauthorized).Methods("GET")
}

func (c *AccountControllers) Home(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	data := struct {
		IsLoggedIn bool
		Username   string
	}{
		IsLoggedIn: cookie != nil && cookie.Value != "",
		Username:   cookie.Value,
	}
	c.template.ExecuteTemplate(w, "home", data)
}

func (c *AccountControllers) CreateForm(w http.ResponseWriter, r *http.Request) {
	c.template.ExecuteTemplate(w, "register", nil)
}

func (c *AccountControllers) LoginForm(w http.ResponseWriter, r *http.Request) {
	c.template.ExecuteTemplate(w, "login", nil)
}

func (c *AccountControllers) Create(w http.ResponseWriter, r *http.Request) {

	hashedpassword := hash.HashPassword(r.PostFormValue("password"))

	newAccount := models.Account{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: hashedpassword,
	}

	_, AccountErr := c.service.Create(newAccount)
	if AccountErr != nil {
		http.Error(w, AccountErr.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/login"), http.StatusMovedPermanently)
}

func (c *AccountControllers) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	account, err := c.service.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sessionCookie := &http.Cookie{
		Name:     "session",
		Value:    account.Username,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, sessionCookie)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (c *AccountControllers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/unauthorized", http.StatusTemporaryRedirect)
			return
		}

		if cookie.Value == "" {
			http.Redirect(w, r, "/unauthorized", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (c *AccountControllers) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	c.template.ExecuteTemplate(w, "404", nil)
}

func (c *AccountControllers) Unauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	c.template.ExecuteTemplate(w, "401", nil)
}
