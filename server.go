// File: server.go
package main

import (
	"log"
	"net/http"
)

// Start server and forward requests to other handlers.
func handleRequests() {
	log.Println("Starting server on :8888...")
	http.HandleFunc("/user/login", loginHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func main() {
	handleRequests()
}