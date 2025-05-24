package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var errorRateLimitExceeded = errors.New("rate limit for GitHub API was likely exceeded")

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		os.Stdout.WriteString("Usage: gcdstats [user] [repository, optional]\n")
		os.Exit(1)
	}

	var repos []repository

	if len(args) == 2 {
		repos = []repository{{Name: args[0] + "/" + args[1]}}
	} else if strings.Contains(args[0], "/") {
		repos = []repository{{Name: args[0]}}
	} else {
		allRepos, err := fetchRepositories(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch repositories: %v\n", err)
			os.Exit(1)
		}

		repos = allRepos
	}

	err := fetchStatistics(repos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch statistics: %v\n", err)
		os.Exit(1)
	}
}
