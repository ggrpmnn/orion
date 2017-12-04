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
	router.HandleFunc("/analyze", analyzeHandler).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)

	// for serving static files (HTML, CSS, etc.)
	router.Handle("/css/{file}", http.FileServer(http.Dir("")))
	router.Handle("/img/{file}", http.FileServer(http.Dir("")))

	log.Print("listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
