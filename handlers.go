package main

import (
	"fmt"
	"log"
	"net/http"
)

// index handles the requests to the main page
func index(w http.ResponseWriter, r *http.Request) {
	log.Print("received homepage request from " + r.RemoteAddr)
	fmt.Fprintln(w, "Orion Homepage")
}

// notFound handles any requests to unknown resources
func notFound(w http.ResponseWriter, r *http.Request) {
	log.Print("received unknown " + r.Method + " request for '" + r.URL.String() + "' from " + r.RemoteAddr)
	sendResponse(w, `{"error": "requested resource does not exist"}`, http.StatusNotFound)
}

// analyze responds immediately and begins the code analysis
// for the specified repo
func analyze(w http.ResponseWriter, r *http.Request) {
	log.Print("received analyis request")
	sendResponse(w, `{"message": "received request to analyze code"}`, http.StatusOK)

	// get repo data from Github

	// determine what code the repo has (can be mutliple types)

	// for each language
}
