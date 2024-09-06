package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 注意: 本番環境では適切に制限してください
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to upgrade to WebSocket")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	log.Println("WebSocket connection established")
	defer conn.Close()

	for {
		time.Sleep(5 * time.Second)
		log.Println("Sending message")
		message := `
		   <div>
     <p class="text-sm mb-1">New message from test topic: Hello, World!</p>
   </div>`
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}
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
