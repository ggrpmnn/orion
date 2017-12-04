package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/template"
	sj "github.com/bitly/go-simplejson"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("html/*"))
	log.Println("HTML templates parsed")
}

// index handles the requests to the main page
func index(w http.ResponseWriter, r *http.Request) {
	log.Printf("received homepage request from " + r.RemoteAddr)
	//fmt.Fprintln(w, "Orion Homepage")
	templates.ExecuteTemplate(w, "header", nil)
	templates.ExecuteTemplate(w, "index", nil)
	templates.ExecuteTemplate(w, "footer", nil)
}

// notFound handles any requests to unknown resources
func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("received unknown " + r.Method + " request for '" + r.URL.String() + "' from " + r.RemoteAddr)
	sendResponse(w, `{"error": "requested resource does not exist"}`, http.StatusNotFound)
}

// analyze responds immediately and begins the code analysis
// for the specified repo
func analyze(w http.ResponseWriter, r *http.Request) {
	log.Printf("received analyis request from " + r.RemoteAddr)

	// parse the request body and convert to JSON
	if r.Body == nil {
		log.Printf("received nil request body")
		sendResponse(w, `{"error": "received empty request body"}`, http.StatusBadRequest)
		return
	}
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("couldn't read byte data from request")
		sendResponse(w, `{"error": "failed to parse request data; try resending the message"}`, http.StatusInternalServerError)
		return
	}
	r.Body.Close()
	js, err := sj.NewJson(bytes)
	if err != nil {
		log.Printf("couldn't marshal request byte data to JSON object")
		sendResponse(w, `{"error": "failed to convert request data to JSON; try resending the message"}`, http.StatusInternalServerError)
		return
	}

	// check the event type; we only care about newly opened PRs (to avoid double commenting)
	if js.Get("action").MustString() != "opened" {
		log.Printf("received a message for an event that was not 'opened'; skipping")
		sendResponse(w, `{"message": "event received was not for PR open; skipping this event"}`, http.StatusOK)
		return
	}

	// begin analysis of repo code
	sendResponse(w, `{"message": "received request to analyze code; beginning analysis - findings will be posted to a comment on the PR"}`, http.StatusOK)
	analyzeCode(js)
}

// formResponse creates and writes the HTTP response message
func sendResponse(w http.ResponseWriter, m string, c int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	w.Write([]byte(m + "\n"))
	return
}
