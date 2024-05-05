package main

import (
	"log"
	"net/http"

	"github.com/ms-scalex-assignment/handler"
)

func main() {
	log.Println("Welcome to ScaleX eBook Library...")

	// initializing routes
	handler.Routes()

	log.Println("Starting server at http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
