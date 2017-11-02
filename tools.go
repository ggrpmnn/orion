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

func init() {
	var err error

	err = exec.Command("/usr/bin/which", "gas").Run()
	if err != nil {
		log.Fatalf("gas (Go(lang) source tool) not installed; exiting")
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
	rx := regexp.MustCompile(`\[([\/\w\.]+):(\d+)\] - (.*)`)
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
