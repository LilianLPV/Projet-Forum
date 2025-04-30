package main

import (
	"fmt"
	"net/http"
	"Projet-Forum/backend/handlers" // Import du package handlers
)

func main() {
	http.HandleFunc("/", handlers.AccueilHandler)
	fmt.Println("Le serveur est lancé http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
