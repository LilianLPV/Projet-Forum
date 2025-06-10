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
	r.HandleFunc("/signup", c.CreateForm).Methods("GET")
	r.HandleFunc("/signup/treatment", c.Create).Methods("POST")
	r.HandleFunc("/login", c.LoginForm).Methods("GET")
	r.HandleFunc("/login", c.Login).Methods("POST")
}

func (c *AccountControllers) CreateForm(w http.ResponseWriter, r *http.Request) {
	c.template.ExecuteTemplate(w, "account.create", nil)
}
func (c *AccountControllers) LoginForm(w http.ResponseWriter, r *http.Request) {
	c.template.ExecuteTemplate(w, "account.login", nil)
}

func (c *AccountControllers) Create(w http.ResponseWriter, r *http.Request) {

	hashedpassword, passwordErr := hash.HashPassword(r.PostFormValue("password"))
	if passwordErr != "" {
		http.Error(w, "Erreur - Le mot de passe est invalide", http.StatusBadRequest)
		return
	}
	newAccount := models.Account{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: hashedpassword,
	}

	productId, productErr := c.service.Create(newAccount)
	if productErr != nil {
		http.Error(w, productErr.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(fmt.Sprintf("/product/%d", productId))

	http.Redirect(w, r, fmt.Sprintf("/product/%d", productId), http.StatusMovedPermanently)
}
func (c *AccountControllers) Login(w http.ResponseWriter, r *http.Request) {

	r.FormValue("mail")
}
