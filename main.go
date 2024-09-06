package main

import (
	"embed"
	"fmt"
	"html/template"
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

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("homeHandler")
	// Load the template
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get query parameter
	name := r.URL.Query().Get("name")

	// Render the template
	err = tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Name": name,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("aboutHandler")
	// Load the template
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	err = tmpl.ExecuteTemplate(w, "about.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("contactHandler")
	// Load the template
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	err = tmpl.ExecuteTemplate(w, "contact.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
