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

	accountService := services.InitAccountService(db)
	accountController := controllers.InitAccountController(accountService, temp)

	feedService := services.InitFeedService(db)
	feedController := controllers.InitFeedController(feedService, temp)

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		temp.ExecuteTemplate(w, "404", nil)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		temp.ExecuteTemplate(w, "401", nil)
	})

	accountController.AccountRouter(router)
	feedController.FeedRouter(router)

	fs := http.FileServer(http.Dir("static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	serveErr := http.ListenAndServe(":8080", router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur - %s", serveErr.Error())
	}
	fmt.Println("Serveur lancé : http://localhost:8080")
}
