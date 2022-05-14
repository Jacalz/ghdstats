package main

import (
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		os.Stdout.WriteString("Usage: gcdstats [user] [repository, optional]\n")
		return
	}

	var repos []repository

	if len(args) == 2 {
		repos = []repository{{Name: args[0] + "/" + args[1]}}
	} else if strings.Contains(args[0], "/") {
		repos = []repository{{Name: args[0]}}
	} else {
		repos = fetchRepositories(args[0])
	}

	fetchStatistics(repos)
}
