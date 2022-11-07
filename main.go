package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jacob4649/cactus-of-things-backend/cactus-of-things-backend/sensor"
)

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)

	datastore := sensor.SensorDatastore{ProjectID: "cactus-of-things-backend"}
	readingGetter, readingSetter := sensor.GetHandlers(&datastore)

	http.HandleFunc("/readings/store", readingSetter)
	http.HandleFunc("/readings/{start}/to/{end}", readingGetter)

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

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s From Go!\n", name)
}
