package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const apiEndpoint = "https://api.github.com/users/"

type repository struct {
	Name string `json:"full_name"`
}

func fetchRepositories(user string) []repository {
	resp, err := http.Get(apiEndpoint + user + "/repos")
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

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		fmt.Println("Usage: gcdstats [user] [repository, optional]")
		return
	}

	repos := []repository{}

	if len(args) == 2 {
		repos = append(repos, repository{apiEndpoint + args[0] + "/" + args[1]})
	} else {
		repos = append(repos, fetchRepositories(args[0])...)
	}

	fmt.Println(repos)
}
