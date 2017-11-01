package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// Finding represents a result from one of the analysis tools
type Finding struct {
	// the filename of the finding
	File string
	// the line number of the finding
	Line string
	// the text of the finding
	Text string
}

// AnalysisTool represents a single tool for performing code analysis
type AnalysisTool struct {
	// the language that the tool analyzes
	Language string
	// the name of the tool
	Name string
	// the path (absolute) to the tool on this system
	Path string
	// the arguments to be passed to the tool on the command line
	Args []string
	// the output of the tool
	Output []Finding
}

var tools map[string]AnalysisTool

func init() {
	var err error
	tools = make(map[string]AnalysisTool)

	configBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("failed to load ./config.json file")
	}
	toolsList := make([]AnalysisTool, 0)
	err = json.Unmarshal(configBytes, &toolsList)
	if err != nil {
		log.Fatal("failed to marshal config file; check the file and try again")
	}
	// convert to map for faster lookup later
	for _, tool := range toolsList {
		tools[tool.Language] = tool
	}
	for _, tool := range tools {
		err = exec.Command(tool.Path, tool.Args...).Run()
		if err != nil {
			log.Fatalf("specified tool '%s' is not installed on this system; install the tool or remove it from the configuration", tool.Name)
		}
	}
}

// analyzeGo utilizes the GoAST package to analyze Go(lang) code
func analyzeGo() ([]Finding, error) {
	cmd := exec.Command("gas", "-skip=tests*", "-fmt=json", "./...")
	resBytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	resStr := string(resBytes)

	// parse results
	findings := make([]Finding, 1)
	// output line format/example: [/path/to/file:123] - Errors unhandled. (Confidence: HIGH, Severity: LOW)
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
