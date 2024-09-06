package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/websocket"
	"google.golang.org/api/iterator"
)

var (
	topicList     map[string]string = make(map[string]string)
	topic         string
	hasSubscribed bool = false
)

func messageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("homeHandler")
	// Load the template
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	client, err := pubsub.NewClient(ctx, ProjectID())
	if err != nil {
		log.Printf("Error creating Pub/Sub client: %v", err)
		return
	}
	defer client.Close()

	topics := make([]string, 0)

	it := client.Subscriptions(ctx)
	for {
		log.Println("topicList")
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error getting subscription: %v", err)
			return
		}

		log.Printf("sub.ID(): %s", sub.ID())

		if strings.HasPrefix(sub.ID(), "debug-") {
			topic := strings.Replace(sub.ID(), "debug-", "", 1)
			topics = append(topics, topic)

			topicList[topic] = sub.ID()
		}

		log.Printf("topicList: %+v", topicList)
	}

	log.Printf("topics: %+v", topics)

	// Render the template
	err = tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Topics": topics,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("subscribeHandler")
	topic = r.FormValue("topic")
	log.Printf("Topic: %s", topic)

	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully subscribed to topic: " + topic))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
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

	ctx := r.Context()
	client, err := pubsub.NewClient(ctx, ProjectID())
	if err != nil {
		log.Printf("Error creating Pub/Sub client: %v", err)
		return
	}
	defer client.Close()

	for {
		if topic == "" {
			log.Println("Waiting for topic...")
			time.Sleep(5 * time.Second)
			continue
		}
		if hasSubscribed {
			return
		}
		break
	}

	sub := client.Subscription(fmt.Sprintf("debug-%s", topic))

	// メッセージ受信用のチャネル
	messages := make(chan *pubsub.Message)

	// Pub/Subからメッセージを受信するgoroutine
	go func() {
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			messages <- msg
			msg.Ack()
		})
		if err != nil {
			log.Printf("Error receiving messages: %v", err)
		}
	}()

	for {
		select {
		case msg := <-messages:
			log.Printf("Received message: %+v", msg)

			// JSONをインデント
			var prettyJSON bytes.Buffer
			err := json.Indent(&prettyJSON, msg.Data, "", "  ")
			if err != nil {
				log.Printf("Error indenting JSON: %v", err)
				continue
			}

			// HTMLエスケープしてXSS攻撃を防ぐ
			escapedJSON := html.EscapeString(prettyJSON.String())
			publishTime := msg.PublishTime

			// メッセージをHTMLに整形
			htmlMessage := fmt.Sprintf(`
						<div id="messages" hx-swap-oob="beforeend">
								<div class="bg-white p-4 rounded-md shadow-sm border border-gray-200 overflow-x-auto">
										<p class="text-xs text-gray-500 mb-2">%s</p>
										<pre class="text-sm text-gray-800 whitespace-pre-wrap"><code>%s</code></pre>
								</div>
						</div>
				`, publishTime.Format("2006-01-02 15:04:05"), escapedJSON)

			// WebSocketクライアントにメッセージを送信
			err = conn.WriteMessage(websocket.TextMessage, []byte(htmlMessage))
			if err != nil {
				log.Printf("Error sending message to WebSocket: %v", err)
				return
			}
		case <-r.Context().Done():
			log.Println("WebSocket connection closed")
			return
		}
	}
}
