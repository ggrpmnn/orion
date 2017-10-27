package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	sj "github.com/bitly/go-simplejson"
	//"github.com/GoASTScanner/gas"
)

// Finding represents a result from one of the analysis tools
type Finding struct {
	File string
	Line string
	Text string
}

// ensure tools are installed before running any analysis ops
func init() {
	// git
	err := exec.Command("which", "git").Run()
	if err != nil {
		log.Fatal("error: git not installed")
	}
	// GoAST (gas)
	err = exec.Command("which", "gas").Run()
	if err != nil {
		log.Fatal("error: gas not installed")
	}
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
	os.Mkdir(workDir, 0700)
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
	analysisFindings := make(map[string][]Finding)
	for language := range languages {
		switch language {
		case "Go":
			log.Printf("%s - running Go(lang) code analysis", repoName)
			findings, err := runGoAnalysis()
			if err != nil {
				log.Printf("%s - error while running Go gas tool: %s", repoName, err)
			}
			analysisFindings["Go"] = findings
		default:
			log.Printf("%s - in language loop: in language default with value %s", repoName, language)
		}
	}

	if analysisFindings != nil {
		for lang, af := range analysisFindings {
			fmt.Printf("%s:\n", lang)
			for _, f := range af {
				fmt.Printf("%v\n", f)
			}
		}
	}

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
func runGoAnalysis() ([]Finding, error) {
	cmd := exec.Command("gas", "-skip=tests*", "-fmt=json", "./...")
	resBytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	resStr := string(resBytes)

	// parse results
	findings := make([]Finding, 1)
	// output line format: [/path/to/file:123] - Errors unhandled. (Confidence: HIGH, Severity: LOW)
	rx := regexp.MustCompile(`\[([\/\w\.]+):(\d+)\] - (.*) \(.*\)`)
	scan := bufio.NewScanner(strings.NewReader(resStr))
	for scan.Scan() {
		if scan.Text() == "\n" {
			continue
		}
		matches := rx.FindStringSubmatch(scan.Text())
		findings = append(findings, Finding{File: matches[1], Line: matches[2], Text: matches[3]})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	return findings, nil
}
