package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

//go:embed templates
var templateFS embed.FS

func main() {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/about", aboutHandler)
	router.HandleFunc("/contact", contactHandler)
	router.HandleFunc("/subscribe", subscribeHandler).Methods("POST") // Add subscribe route

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
