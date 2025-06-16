package controllers

import (
	"encoding/json"
	"forum/auth"
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
	r.HandleFunc("/", c.Home).Methods("GET")
	r.HandleFunc("/signup", c.CreateForm).Methods("GET")
	r.HandleFunc("/login", c.LoginForm).Methods("GET")
	r.HandleFunc("/signup/treatment", c.Create).Methods("POST")
	r.HandleFunc("/login/treatment", c.Login).Methods("POST")
	r.HandleFunc("/logout", c.Logout).Methods("POST")
	r.HandleFunc("/profile", c.ProfilePage).Methods("GET")

}

func (c *AccountControllers) ProfilePage(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	account, err := c.service.ReadById(claims.Sub)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du profil", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username       string
		Email          string
		Role           string
		ProfilePicture string
		Bio            string
	}{
		Username:       account.Username,
		Email:          account.Email,
		Role:           account.Role,
		ProfilePicture: account.Profilepicture,
		Bio:            account.Bio,
	}

	c.template.ExecuteTemplate(w, "profile", data)
}

func (c *AccountControllers) Home(w http.ResponseWriter, r *http.Request) {
	userData := struct {
		IsLoggedIn bool
		Username   string
		Role       string
	}{}

	token, err := r.Cookie("token")
	if err == nil && token.Value != "" {
		claims, err := auth.ValidateToken(token.Value)
		if err == nil {
			userData.IsLoggedIn = true
			userData.Role = claims.Role
			account, err := c.service.ReadById(claims.Sub)
			if err == nil {
				userData.Username = account.Username
			}
		}
	}

	c.template.ExecuteTemplate(w, "home", userData)
}

func (c *AccountControllers) CreateForm(w http.ResponseWriter, r *http.Request) {
	c.template.ExecuteTemplate(w, "signup", nil)
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
		Role:     "user",
	}

	_, AccountErr := c.service.Create(newAccount)
	if AccountErr != nil {
		http.Error(w, AccountErr.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func (c *AccountControllers) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	account, err := c.service.Login(email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(account.Id, account.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	response := struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}{
		Username: account.Username,
		Role:     account.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *AccountControllers) CreatePost(w http.ResponseWriter, r *http.Request) {
}

func (c *AccountControllers) GetPosts(w http.ResponseWriter, r *http.Request) {
}

func (c *AccountControllers) GetPost(w http.ResponseWriter, r *http.Request) {
}

func (c *AccountControllers) UpdatePost(w http.ResponseWriter, r *http.Request) {
}

func (c *AccountControllers) DeletePost(w http.ResponseWriter, r *http.Request) {
}

func (c *AccountControllers) Profile(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	account, err := c.service.ReadById(claims.Sub)
	if err != nil {
		http.Error(w, "Error retrieving profile", http.StatusInternalServerError)
		return
	}

	response := struct {
		Username       string
		Email          string
		Role           string
		ProfilePicture string
		Bio            string
	}{
		Username:       account.Username,
		Email:          account.Email,
		Role:           account.Role,
		ProfilePicture: account.Profilepicture,
		Bio:            account.Bio,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *AccountControllers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	var updateData struct {
		Bio string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	err = c.service.UpdateProfile(claims.Sub, updateData.Bio)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour du profil", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profil mis à jour avec succès"})
}

func (c *AccountControllers) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (c *AccountControllers) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	c.template.ExecuteTemplate(w, "404", nil)
}

func (c *AccountControllers) Unauthorized(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	c.template.ExecuteTemplate(w, "401", nil)
}
