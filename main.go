package main

import (
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		os.Stdout.WriteString("Usage: gcdstats [user] [repository, optional]")
		return
	}

	var repos []repository

	if len(args) == 2 {
		repos = []repository{{Name: args[0] + "/" + args[1]}}
	} else {
		repos = fetchRepositories(args[0])
	}

	fetchStatistics(repos)
}
