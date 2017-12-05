package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	sj "github.com/bitly/go-simplejson"
	"github.com/ggrpmnn/orion/src/github"
)

var orionWorkspace = os.Getenv("ORION_WORKSPACE")

func init() {
	if orionWorkspace == "" {
		orionWorkspace = "."
	}

	err := exec.Command("/usr/bin/which", "git").Run()
	if err != nil {
		log.Fatal("error: git not installed or not added to PATH")
	}
}

// analyzeCode is called when the API receives a webhook message to the /analyze
// endpoint; for actual analysis functionality, see language-specific functions in tools.go
func analyzeCode(json *sj.Json) {
	repoName := github.GetRepoName(json)
	log.Printf("%s - beginning overall code analysis", repoName)

	gitURL := github.GetRepoURL(json)
	if gitURL == "" {
		log.Printf("%s - failed to retrieve git URL from JSON message", repoName)
		return
	}
	// create SHA hash of message and use it as the name of the temporary work directory
	messageHash := sha256.New()
	jsBytes, err := json.Encode()
	if err != nil {
		log.Printf("%s - failed to convert GitHub JSON to byte data for calculating workspace hash", repoName)
		return
	}
	messageHash.Write(jsBytes)
	workDir := fmt.Sprintf("%s/%x", orionWorkspace, messageHash.Sum(nil))
	os.Mkdir(workDir, 0700)
	os.Chdir(workDir)
	defer cleanup(workDir)
	err = exec.Command("git", "clone", addGitHubCredsToURL(gitURL)).Run()
	if err != nil {
		log.Printf("%s - failed to clone repository from target URL '%s': %s", repoName, gitURL, err)
		return
	}

	languages, err := github.GetLanguageMapping(json)
	if err != nil {
		log.Printf("%s - failed to retrieve repository code composition: %s", repoName, err)
		return
	}
	analysisFindings := make(map[string][]Finding)
	for language := range languages {
		switch language {
		case "Go":
			log.Printf("%s - running Go(lang) analysis", repoName)
			findings, err := analyzeGo(repoName)
			if err != nil {
				log.Printf("%s - error while running Go gas tool: %s", repoName, err)
			}
			if len(findings) > 0 {
				analysisFindings["Go"] = findings
			}
		default:
			log.Printf("%s - in language loop: in language default; found '%s'", repoName, language)
		}
	}

	if len(analysisFindings) == 0 {
		log.Printf("%s - no results returned from scan(s); PR seems clean", repoName)
	} else {
		err = postComment(json, analysisFindings, repoName)
		if err != nil {
			log.Printf(err.Error())
		}
	}

	log.Printf("%s - finishing overall code analysis", repoName)
}

// addGitHubCredsToURL takes the GitHub URL and adds credentials so that it takes the
// form `http(s)://username:password@github.com/...`; this is used to avoid a state
// where credentials are requested when attempting to clone a repo, which could cause
// the git process to wait for user input and stall
func addGitHubCredsToURL(url string) string {
	pieces := strings.Split(url, "://")
	pieces[1] = fmt.Sprintf("%s:%s@%s", github.GitHubUser, github.GitHubToken, pieces[1])
	return strings.Join(pieces, "://")
}

// composeCommentText creates the content of a comment message using the findings from the analysis
func composeCommentText(findings map[string][]Finding) string {
	body := "Hi, I'm Orion, a code-analysis application. When you registered your pull request, your code was scanned and the following issues were found.\\n\\n"
	for language, findings := range findings {
		body += fmt.Sprintf("%s:\\n", language)
		for _, finding := range findings {
			body += fmt.Sprintf("* `%s`, line %s: %s\\n", finding.File, finding.Line, finding.Text)
		}
		body += "\\n\\n"
	}
	body += "Please fix these issues before merging this PR, if possible. Thanks!"

	return body
}

// postComment gets the repo's comment URL from the message JSON, composes the comment body,
// then posts the comment (with any findings) to the PR
func postComment(json *sj.Json, findings map[string][]Finding, repoName string) error {
	commentsURL := json.Get("pull_request").Get("comments_url").MustString()
	if commentsURL == "" {
		return fmt.Errorf("%s - failed to retrieve comments URL from JSON message", repoName)
	}
	body := `{"body": "` + composeCommentText(findings) + `"}`
	req, err := http.NewRequest("POST", commentsURL, strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s - failed to create POST request", repoName)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+github.GitHubToken)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s - received error when  POSTing comment: %s", repoName, err)
	}
	if res.StatusCode > 299 {
		if res == nil {
			err = fmt.Errorf("%s - received error status (%d) when POSTing comment", repoName, res.StatusCode)
		} else {
			bytes, _ := ioutil.ReadAll(res.Body)
			return fmt.Errorf("%s - response error: %s", repoName, bytes)
		}
	}
	return nil
}

// cleanup removes new folders and cloned code when the analysis is done; this function preferably
// should be 'defer'ed to ensure that it is run even when a panic/unexpected error occurs
func cleanup(dir string) {
	os.Chdir("../")
	os.RemoveAll(dir)
}
