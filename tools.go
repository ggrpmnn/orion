package main

import (
	"bufio"
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

// Tool represents a single tool for performing code analysis
type Tool struct {
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

var gasCmdPath string

func init() {
	var err error

	outBytes, err := exec.Command("/usr/bin/which", "gas").Output()
	if err != nil {
		log.Fatalf("gas (Go(lang) source tool) not installed; exiting")
	}
	gasCmdPath = string(outBytes)
}

// analyzeGo utilizes the GoAST package to analyze Go(lang) code
func analyzeGo() ([]Finding, error) {
	cmd := exec.Command(gasCmdPath, "-skip=tests*", "-fmt=json", "./...")
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
