package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"meMessages/middlewares"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
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
