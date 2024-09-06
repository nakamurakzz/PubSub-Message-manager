package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("homeHandler")
	// Load the template
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get topic list (replace with your actual logic)
	topics := []string{"Topic1", "Topic2", "Topic3", "Topic4"}

	// Render the template
	err = tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Topics": topics,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type Subscriber struct {
	Topic string `json:"topic"`
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("subscribeHandler")

	// フォームデータからトピックを取得
	topic := r.FormValue("topic")
	if topic == "" {
		log.Println("Error: No topic provided")
		http.Error(w, "No topic provided", http.StatusBadRequest)
		return
	}

	log.Printf("Received topic: %s\n", topic)

	// Replace with your actual logic to register the subscriber
	fmt.Printf("Subscriber registered for topic: %s\n", topic)

	// Send a success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Subscriber registered successfully"))
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
