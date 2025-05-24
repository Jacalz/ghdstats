package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type repository struct {
	Name string `json:"full_name"`
}

func fetchRepositories(user string) []repository {
	const usersApiEndpoint = "https://api.github.com/users/"
	resp, err := http.Get(usersApiEndpoint + user + "/repos")
	if err != nil {
		os.Stderr.WriteString("Error: The HTTP get request failed. Error message: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	defer resp.Body.Close()
	repos := []repository{}

	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		os.Stderr.WriteString("Error: Failed to decode JSON data. A likely culprit is that the GitHub API limit was likely reached.\nTry again in a few hours. Error message: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	return repos
}
