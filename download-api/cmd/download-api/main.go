package main

import (
	"download-api/internal/config"
	"download-api/internal/handlers"
	"log"
	"net/http"
)

func main() {
	config.LoadConfig()

	host := config.K.String("HOST")
	if host == "" {
		host = "0.0.0.0" // Default host
	}

	port := config.K.String("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	address := host + ":" + port

	http.HandleFunc("/download", handlers.DownloadZipHandler)
	http.HandleFunc("/livez", livezHandler)
	http.HandleFunc("/readyz", readyzHandler)

	log.Printf("Starting server on %s\n", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func livezHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
