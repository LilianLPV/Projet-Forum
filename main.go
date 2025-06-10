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

	temp, tempErr := template.ParseGlob("./templates/*.html")
	if tempErr != nil {
		log.Fatalf("Erreur chargement serveur - %s", tempErr.Error())
	}

	AccountService := services.InitAccountService(db)

	accountController := controllers.InitAccountController(AccountService, temp)

	router := mux.NewRouter()
	accountController.AccountRouter(router)

	// Lancement du serveur
	serveErr := http.ListenAndServe(":8080", router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur - %s", serveErr.Error())
	}
	fmt.Println("Serveur lancé : http://localhost:8080")
}


}
