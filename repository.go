package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const usersApiEndpoint = "https://api.github.com/users/"

type repository struct {
	Name string `json:"full_name"`
}

func fetchRepositories(user string) []repository {
	resp, err := http.Get(usersApiEndpoint + user + "/repos")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: The HTTP get request failed. Error message: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	repos := []repository{}

	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to decode JSON data. A likely culprit is that the GitHub API limit was reached.\nTry again in a few hours. Error message: %v\n")
		os.Exit(1)
	}

	return repos
}
