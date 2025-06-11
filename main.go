package main

import (
	"fmt"
	"forum/config"
	"forum/controllers"
	"forum/services"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

func main() {
	config.LoadEnv()

	db, dbErr := config.InitDB()
	if dbErr != nil {
		log.Fatal(dbErr.Error())
		return
	}
	defer db.Close()

	temp, tempErr := template.ParseGlob("./template/*")
	if tempErr != nil {
		log.Fatalf("Erreur chargement serveur - %s", tempErr.Error())
	}

	AccountService := services.InitAccountService(db)
	accountController := controllers.InitAccountController(AccountService, temp)

	router := mux.NewRouter()

	// Configuration des gestionnaires d'erreurs personnalisés
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		temp.ExecuteTemplate(w, "404", nil)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		temp.ExecuteTemplate(w, "401", nil)
	})

	// Configuration des routes
	accountController.AccountRouter(router)

	// Configuration des fichiers statiques
	fs := http.FileServer(http.Dir("static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Lancement du serveur
	serveErr := http.ListenAndServe(":8080", router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur - %s", serveErr.Error())
	}
	fmt.Println("Serveur lancé : http://localhost:8080")
}
