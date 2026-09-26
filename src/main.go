package main

import (
	"log"
	"meMessages/db"
	"meMessages/middlewares"
	"net/http"
	"os"
	"meMessages/hub"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //  this is for dev, we gotta tighten this in prod (Which i will never up to)
	},
}

func wsHandler(hub *hub.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")
		if userId == "" {
			http.Error(w, "user_id is required", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v", err)
			return 
		}
		defer conn.Close()

		hub.Register(userId, conn)
		defer hub.Unregister(userId)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break // connection refused
			}

			log.Printf("received from %s -> %s", userId, msg)
		}
	}
}

func main() {
	hosts := os.Getenv("CQL_HOSTS")
	keyspace := os.Getenv("CQL_KEYSPACE")

	if hosts == "" || keyspace == "" {
		log.Fatal("CQL_HOSTS and CQL_KEYSPACE MUST be set")
		return
	}
	
	db.EnsureKeyspace(hosts, keyspace)

	session := db.Connect(hosts, keyspace)
	defer session.Close()

	db.Migrate(session)
	
	// ---- migration over ----
	
	// new ws hub ->
	h := hub.NewHub()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ws", wsHandler(h))
	
	handler := middlewares.LoggerMiddleware(mux)

	server := http.Server{
		Addr: ":4444",
		Handler: handler,
	}
	
	log.Printf("Server is started on :4444")
	server.ListenAndServe()
}
