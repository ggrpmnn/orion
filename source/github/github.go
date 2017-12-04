package github

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// GetLanguages queries GitHub for the repository's language composition; see
// https://developer.github.com/v3/repos/#list-languages for more information on the service
func GetLanguages(endpoint string) (map[string]int, error) {
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
