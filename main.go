package main

import (
	"fmt"
	"forum/handlers" // Assurez-vous que le chemin d'importation correspond à votre structure de projet
	"net/http"
)

func main() {
	
	http.HandleFunc("/", handlers.AccueilHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Le serveur est lancé http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
