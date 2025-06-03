package main

import (
	"fmt"
	"forum/handlers"
	"net/http"
)

func main() {

	http.HandleFunc("/", handlers.AccueilHandler)
	fmt.Println("Le serveur est lancé http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
