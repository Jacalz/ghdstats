package main

import (
	"encoding/json"
	"net/http"
)

type repository struct {
	Name string `json:"full_name"`
}

func fetchRepositories(user string) []repository {
	const usersApiEndpoint = "https://api.github.com/users/"
	resp, err := http.Get(usersApiEndpoint + user + "/repos")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	repos := []repository{}

	err = json.NewDecoder(resp.Body).Decode(&repos)
	if err != nil {
		panic(err)
	}

	return repos
}
