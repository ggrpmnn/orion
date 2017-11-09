package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = "8080"

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)
	router.HandleFunc("/analyze", analyze).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)

	log.Print("listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
