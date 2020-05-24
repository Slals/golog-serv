package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/logs", logHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	log.Printf("Listening to :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
