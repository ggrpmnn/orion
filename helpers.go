package main

import "net/http"

// formResponse creates and writes the HTTP response message
func sendResponse(w http.ResponseWriter, m string, c int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	w.Write([]byte(m + "\n"))
	return
}
