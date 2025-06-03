package main

import (
	"fmt"
	"forum/handlers" // Assurez-vous que le chemin d'importation correspond à votre structure de projet
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	
	http.HandleFunc("/", handlers.AccueilHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Le serveur est lancé http://localhost:8080")
	http.ListenAndServe(":8080", nil)

	router := mux.NewRouter()
	serveErr := http.ListenAndServe(":8080", router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur - %s", serveErr.Error())
	}
	fmt.Println("Serveur lancé : http://localhost:8080")
}
