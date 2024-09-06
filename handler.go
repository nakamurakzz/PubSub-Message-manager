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
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/websocket"
	"google.golang.org/api/iterator"
)

var (
	topicList     map[string]string = make(map[string]string)
	subCancel     context.CancelFunc
	subMutex      sync.Mutex
	wsConnections []*websocket.Conn
	wsConnMutex   sync.Mutex
)

func messageHandler(w http.ResponseWriter, r *http.Request) {
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

	err = tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{
		"Topics": topics,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	newTopic := r.FormValue("topic")
	log.Printf("New Topic: %s", newTopic)

	subMutex.Lock()
	defer subMutex.Unlock()

	if subCancel != nil {
		subCancel()
		subCancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	subCancel = cancel

	go subscribeToTopic(ctx, newTopic)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully subscribed to topic: " + newTopic))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	wsConnMutex.Lock()
	wsConnections = append(wsConnections, conn)
	wsConnMutex.Unlock()

	defer func() {
		wsConnMutex.Lock()
		for i, c := range wsConnections {
			if c == conn {
				wsConnections = append(wsConnections[:i], wsConnections[i+1:]...)
				break
			}
		}
		wsConnMutex.Unlock()
	}()

	// WebSocket接続が閉じられるまで待機
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func subscribeToTopic(ctx context.Context, topic string) {
	client, err := pubsub.NewClient(ctx, ProjectID())
	if err != nil {
		log.Printf("Error creating Pub/Sub client: %v", err)
		return
	}
	defer client.Close()

	sub := client.Subscription(fmt.Sprintf("debug-%s", topic))

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		sendMessageToWebSocket(msg)
		msg.Ack()
	})

	if err != nil && err != context.Canceled {
		log.Printf("Error receiving messages: %v", err)
	}
}

func sendMessageToWebSocket(msg *pubsub.Message) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, msg.Data, "", "  ")
	if err != nil {
		log.Printf("Error indenting JSON: %v", err)
		return
	}

	escapedJSON := html.EscapeString(prettyJSON.String())
	publishTime := msg.PublishTime

	htmlMessage := fmt.Sprintf(`
		<div id="messages" hx-swap-oob="beforeend">
			<div class="bg-white p-4 rounded-md shadow-sm border border-gray-200 overflow-x-auto">
				<p class="text-xs text-gray-500 mb-2">%s</p>
				<pre class="text-sm text-gray-800 whitespace-pre-wrap"><code>%s</code></pre>
			</div>
		</div>
	`, publishTime.Format("2006-01-02 15:04:05"), escapedJSON)

	wsConnMutex.Lock()
	defer wsConnMutex.Unlock()

	for _, conn := range wsConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(htmlMessage))
		if err != nil {
			log.Printf("Error sending message to WebSocket: %v", err)
		}
	}
}
