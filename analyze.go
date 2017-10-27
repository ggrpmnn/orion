package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	sj "github.com/bitly/go-simplejson"
	//"github.com/GoASTScanner/gas"
)

// Finding represents a result from one of the analysis tools
type Finding struct {
	File string
	Line string
	Text string
}

// analyzeCode is the function that is called when the API receives a webhook message to the /analyze
// endpoint; for actual analysis functionality, see language-specific functions below
func analyzeCode(json *sj.Json) {
	// get repo data from Github
	repoName := json.Get("repository").Get("name").MustString()
	log.Printf("%s - beginning overall code analysis", repoName)

	htmlURL := json.Get("pull_request").Get("html_url").MustString()
	if htmlURL == "" {
		jsonBytes, _ := json.Encode()
		log.Printf("%s - failed to retrieve repo URL from data: %s", repoName, string(jsonBytes))
		return
	}
	gitURL := json.Get("repository").Get("git_url").MustString()
	if gitURL == "" {
		jsonBytes, _ := json.Encode()
		log.Printf("%s - failed to retrieve git URL from data: %s", repoName, string(jsonBytes))
		return
	}
	// create SHA hash of message and use it as the name of the work directory
	messageHash := sha256.New()
	jsBytes, err := json.Encode()
	if err != nil {
		log.Printf("%s - failed to convert GitHub JSON to byte data for calculating workspace hash", repoName)
		return
	}
	messageHash.Write(jsBytes)
	workDir := fmt.Sprintf("./%x", messageHash.Sum(nil))
	os.Mkdir(workDir, 0755)
	os.Chdir(workDir)
	err = exec.Command("git", "clone", gitURL).Run()
	if err != nil {
		log.Printf("%s - failed to clone repository from target URL '%s': %s", repoName, gitURL, err)
		return
	}

	// determine the code composition and call appropriate analysis functions based on language type
	languages, err := getLanguageComposition(json.Get("repository").Get("languages_url").MustString())
	if err != nil {
		log.Printf("%s - failed to retrieve repository code composition: %s", repoName, err)
		return
	}
	for language := range languages {
		switch language {
		case "Go":
			go runGoAnalysis(json)
		default:
			log.Printf("%s - in language loop: landed in default case with value %s", repoName, language)
		}
	}

	// process Finding outputs, compose comment, send to GitHub URL
	// database recording?

	// cleanup working directory
	os.Chdir("../")
	os.RemoveAll(workDir)

	log.Printf("%s - finishing overall code analysis", repoName)
}

// getLanguageComposition queries GitHub for the repository's language composition; see
// https://developer.github.com/v3/repos/#list-languages for more information on the language service
func getLanguageComposition(endpoint string) (map[string]int, error) {
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

// runGoAnalysis utilizes the GoAST package to analyze Go(lang) code
func runGoAnalysis(json *sj.Json) []Finding {
	repoName := json.Get("repository").Get("name").MustString()
	log.Printf("%s - running Go(lang) code analysis", repoName)

	return nil
}
