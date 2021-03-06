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
	err := exec.Command("/usr/bin/which", "gas").Run()
	if err != nil {
		log.Fatalf("gas (Go(lang) source tool) not installed; exiting")
	}
}

// analyzeGo utilizes the GoAST package to analyze Go(lang) code and returns
// an array of Finding structs containing any results produced by the tool
func analyzeGo(repoName string) ([]Finding, error) {
	cmd := exec.Command("gas", "./...")
	resBytes, _ := cmd.Output()
	resStr := string(resBytes)

	findings := make([]Finding, 0)
	// output line format/example: [/path/to/file:123] - Errors unhandled. (Confidence: HIGH, Severity: LOW)
	rx := regexp.MustCompile(`\[([\S]+):(\d+)\] - (.*)`)
	scan := bufio.NewScanner(strings.NewReader(resStr))
	for scan.Scan() {
		matches := rx.FindStringSubmatch(scan.Text())
		if len(matches) <= 0 {
			continue
		}
		filepath := matches[1]

		findings = append(findings, Finding{File: filepath, Line: matches[2], Text: matches[3]})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}

	return findings, nil
}

// trimFileName looks in the full path for the delim string and then trims the filepath
// accordingly; for example: calling trimFileName("/test/file/path", "file") will return
// "file/path"
func trimFileName(filepath string, delim string) string {
	pathSplit := strings.Split(filepath, "/")
	ok, idx := pathContains(pathSplit, delim)
	if ok {
		return strings.Join(pathSplit[idx+1:], "/")
	}
	return filepath
}

// pathContains returns true if the given path contains the specified value, false otherwise
// used to see if a path (split on the separator) contains a particular folder name
func pathContains(path []string, lookup string) (bool, int) {
	for idx, val := range path {
		if val == lookup {
			return true, idx
		}
	}
	return false, -1
}
