package github

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/bitly/go-simplejson"
)

var (
	// GitHubUser is the GitHub account's username stored in the env
	GitHubUser = os.Getenv("GH_USERNAME")

	// GitHubToken is the GitHub account auth token stored in the env
	GitHubToken = os.Getenv("GH_AUTH_TOKEN")
)

func init() {
	if GitHubUser == "" || GitHubToken == "" {
		log.Fatal("ERROR: GitHub credentials were not supplied to the application")
	}
}

// GetRepoName returns the name of the repository from the webhook message
func GetRepoName(js *simplejson.Json) string {
	return js.Get("repository").Get("name").MustString()
}

// GetRepoURL returns the base URL for the repository from the webhook message;
// this URL specifically is used to clone the repository from GitHub
func GetRepoURL(js *simplejson.Json) string {
	return js.Get("pull_request").Get("head").Get("repo").Get("clone_url").MustString()
}

// GetLanguageMapping queries GitHub for the repository's language composition, and returns
// a map of languages (strings) to their weight in lines of code (int); see
// https://developer.github.com/v3/repos/#list-languages for more information on the service
func GetLanguageMapping(js *simplejson.Json) (map[string]int, error) {
	endpoint := js.Get("repository").Get("languages_url").MustString()
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	languages := make(map[string]int)
	err = json.Unmarshal(body, &languages)
	if err != nil {
		return nil, err
	}
	return languages, nil
}
