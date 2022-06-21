package main

import (
	"dogidentifier/server"
	"dogidentifier/handlers"
	"log"
	"net/http"
	"os"
)


func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/isUrlOfDog", handlers.IsUrlOfDog)

	server := server.NewServer(mux)
	var certificateFile = os.Getenv("IMAGE_IDENTIFIER_CERT_FILE")
	var keyFile = os.Getenv("IMAGE_IDENTIFIER_KEY_FILE")
	log.Println("Starting server...")
	err := server.ListenAndServeTLS(certificateFile, keyFile)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
