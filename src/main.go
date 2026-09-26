package main

import (
	"fmt"
	"log"
	"meMessages/db"
	"meMessages/middlewares"
	"net/http"
	"os"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", greet)
	
	handler := middlewares.LoggerMiddleware(mux)

	server := http.Server{
		Addr: ":4444",
		Handler: handler,
	}
	
	log.Printf("Server is started on :4444")
	server.ListenAndServe()
}
