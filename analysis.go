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
)

func init() {
	var err error

	err = exec.Command("/usr/bin/which", "git").Run()
	if err != nil {
		log.Fatal("error: git not installed or not added to PATH")
	}
}

// analyzeCode is is called when the API receives a webhook message to the /analyze
// endpoint; for actual analysis functionality, see language-specific functions in tools.go
func analyzeCode(json *sj.Json) {
	repoName := json.Get("repository").Get("name").MustString()
	log.Printf("%s - beginning overall code analysis", repoName)

	gitURL := json.Get("pull_request").Get("head").Get("repo").Get("clone_url").MustString()
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
	os.Mkdir(workDir, 0700)
	os.Chdir(workDir)
	defer cleanup(workDir)
	err = exec.Command("git", "clone", gitURL).Run()
	if err != nil {
		log.Printf("%s - failed to clone repository from target URL '%s': %s", repoName, gitURL, err)
		return
	}

	languages, err := getRepoLanguages(json.Get("repository").Get("languages_url").MustString())
	if err != nil {
		log.Printf("%s - failed to retrieve repository code composition: %s", repoName, err)
		return
	}
	analysisFindings := make(map[string][]Finding)
	for language := range languages {
		switch language {
		case "Go":
			log.Printf("%s - running Go(lang) code analysis", repoName)
			findings, err := analyzeGo(repoName)
			if err != nil {
				log.Printf("%s - error while running Go gas tool: %s", repoName, err)
			}
			analysisFindings["Go"] = findings
		default:
			log.Printf("%s - in language loop: in language default; found '%s'", repoName, language)
		}
	}

	// TODO - remove this when comment writing functionality is implemented
	fmt.Println(composeCommentText(analysisFindings))

	log.Printf("%s - finishing overall code analysis", repoName)
}

// cleanup removes new folders and cloned code when the analysis is done; this function preferably
// should be 'defer'ed to ensure that it is run even when a panic/unexpected error occurs
func cleanup(dir string) {
	os.Chdir("../")
	os.RemoveAll(dir)
}

// getRepoLanguages queries GitHub for the repository's language composition; see
// https://developer.github.com/v3/repos/#list-languages for more information on the service
func getRepoLanguages(endpoint string) (map[string]int, error) {
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

func composeCommentText(af map[string][]Finding) string {
	body := "Hi, I'm Orion, a code-analysis application. When you registered your pull request, your code was scanned and the following issues were found.\n\n"
	for language, findings := range af {
		body += fmt.Sprintf("%s:\n", language)
		for _, finding := range findings {
			body += fmt.Sprintf("* %s, line %s: %s\n", finding.File, finding.Line, finding.Text)
		}
		body += "\n\n"
	}

	body += "It's a good idea to fix these issues before merging this PR, if possible. Thanks!"

	return body
}
