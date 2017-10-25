package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/exec"

	sj "github.com/bitly/go-simplejson"
	//"github.com/GoASTScanner/gas"
)

// Finding represents a result from one of the analysis tools
type Finding struct {
	Text string
}

func analyzeCode(json *sj.Json) {
	// 1. get repo data from Github
	// 2. determine what code the repo has (can be mutliple types)
	// 3. for each language, run analysis tool(s) and process output
	// 4. [LATER] consult DB to determine which changes have been reported and drop any that have
	// 5. form output(s) into a comment and post to the PR
	repoName := json.Get("repository").Get("name").MustString()
	log.Printf("%s - beginning code analysis", repoName)

	// pull code repo
	htmlURL := json.Get("pull_request").Get("html_url").MustString()
	if htmlURL == "" {
		jsonBytes, _ := json.Encode()
		log.Printf("%s - unable to retrieve repo URL from data: %s", repoName, string(jsonBytes))
		return
	}
	gitURL := json.Get("repository").Get("git_url").MustString()
	if gitURL == "" {
		jsonBytes, _ := json.Encode()
		log.Printf("%s - unable to retrieve git URL from data: %s", repoName, string(jsonBytes))
		return
	}
	// create SHA hash of message and use it as the name of the work directory
	messageHash := sha256.New()
	jsBytes, err := json.Encode()
	if err != nil {
		log.Printf("%s - could not convert GitHub JSON to byte data for calculating workspace hash", repoName)
		return
	}
	messageHash.Write(jsBytes)
	workDir := fmt.Sprintf("./%x", messageHash.Sum(nil))
	os.Mkdir(workDir, 0755)
	os.Chdir(workDir)
	err = exec.Command("git", "clone", gitURL).Run()
	if err != nil {
		log.Printf("%s - failed to clone repo from target URL '%s': %s", repoName, gitURL, err)
		return
	}

	// begin analysis

	// cleanup working directory
	os.Chdir("../")
	os.RemoveAll(workDir)

	log.Printf("%s - finishing code analysis", repoName)
}
