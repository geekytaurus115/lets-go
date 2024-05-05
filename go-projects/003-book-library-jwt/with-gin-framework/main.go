package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ms-scalex-assignment/with-gin-framework/handler"
)

func main() {

	r := gin.Default()

	// initializing routes
	handler.Routes(r)

	log.Println("Starting server at http://localhost:8080...")
	r.Run(":8080")
}
