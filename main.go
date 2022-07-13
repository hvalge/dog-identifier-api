package main

import (
	"dogidentifier/server"
	"dogidentifier/handlers"
	"log"
	"net/http"
	"os"
);


func main() {
	var certificateFile = os.Getenv("IMAGE_IDENTIFIER_CERT_FILE");
	var keyFile = os.Getenv("IMAGE_IDENTIFIER_KEY_FILE");

	mux := http.NewServeMux();
	mux.HandleFunc("/isUrlOfDog", handlers.IsUrlOfDog);

	server := server.NewServer(mux);
	log.Println("Starting server...");
	
	err := server.ListenAndServeTLS(certificateFile, keyFile);
	if err != nil {
		log.Fatalf("Server failed to start: %v", err);
	}
}
