package main

import (
	"fmt"
	"invento/oauth/logger/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!") // Send a response back to the client
	})

	// Start the server on port 8080
	fmt.Println("Server listening on port 8080")
	go http.ListenAndServe(":8080", nil)
	go handlers.HandleQueueConsumer()
	// Infinite loop
	for {
		select {
		case <-sigChan:
			// Handle the received signal (e.g., cleanup or graceful shutdown)
			log.Println("Received termination signal. Shutting down...")
			// Perform cleanup or graceful shutdown operations here
			log.Println("Shutdown complete.")
			return
		}
	}

}
