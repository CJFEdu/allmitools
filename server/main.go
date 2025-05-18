package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define the port to listen on
	port := 3000

	// Define a basic handler for the root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the AllMiTools server!")
	})

	// Start the server
	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server starting on port %d...\n", port)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
