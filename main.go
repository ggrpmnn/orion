package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Port is the port on which the server will listen
const port = "8080"

func main() {
	// create the router
	router := mux.NewRouter().StrictSlash(true)

	// set up the endpoints
	router.HandleFunc("/", index)
	router.HandleFunc("/analyze", analyze).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)

	// start listening for requests
	log.Print("listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
