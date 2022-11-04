package main

import (
	"log"
	"net/http"
	"github.com/test/transaction/router"
)

func main() {
	// route
	router := router.NewRouter()

	// log server
	log.Printf("Listening on :8000\n")

	// running server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}