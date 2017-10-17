package main

import (
	"log"

	"github.com/bitly/go-simplejson"
	//github.com/GoASTScanner/gas
)

func analyzeCode(json *simplejson.Json) {
	// 1. get repo data from Github
	// 2. determine what code the repo has (can be mutliple types)
	// 3. for each language, run analysis tool(s) and process output
	// 4. form output(s) into a comment and post to the PR
	log.Print("beginning code analysis")

	log.Print("finishing code analysis")
}
