package main

import (
	"encoding/json"
	"net/http"
)

const usersApiEndpoint = "https://api.github.com/users/"

type repository struct {
	Name string `json:"full_name"`
}

func fetchRepositories(user string) ([]repository, error) {
	resp, err := http.Get(usersApiEndpoint + user + "/repos")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	repos := []repository{}

	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}
