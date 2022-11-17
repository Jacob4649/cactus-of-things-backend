package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jacob4649/cactus-of-things-backend/cactus-of-things-backend/sensor"
)

func main() {
	log.Print("starting server...")

	datastore := sensor.SensorDatastore{ProjectID: "cactus-of-things"}
	readingGetter, readingSetter := sensor.GetHandlers(&datastore)

	http.HandleFunc("/readings/store", readingSetter)
	http.HandleFunc("/readings", readingGetter)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
