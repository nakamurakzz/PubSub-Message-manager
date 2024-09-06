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
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./cmd/pubusub-subscriber <project_id> <port(optional)>")
	}
	if len(os.Args) >= 3 {
		os.Setenv("PORT", os.Args[2])
	}
	os.Setenv("PROJECT_ID", os.Args[1])
	router := mux.NewRouter()

	router.HandleFunc("/ws", handleWebSocket)

	// Define routes
	router.HandleFunc("/", messageHandler)
	router.HandleFunc("/subscribe", subscribeHandler).Methods("POST") // Add subscribe route

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func ProjectID() string {
	return os.Getenv("PROJECT_ID")
}

func Port() string {
	return os.Getenv("PORT")
}
