package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type repository struct {
	Name string `json:"full_name"`
}

type statistics struct {
	Assets []struct {
		Name          string    `json:"name"`
		DownloadCount uint      `json:"download_count"`
		Date          time.Time `json:"updated_at"`
	} `json:"assets"`
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

func fetchStatistics(repos []repository) {
	wg := sync.WaitGroup{}
	wg.Add(len(repos))

	const reposApiEndpoint = "https://api.github.com/repos/"
	for _, repo := range repos {
		go func(repourl, reponame string) {
			defer wg.Done()

			fmt.Println(reponame, repourl)

			resp, err := http.Get(repourl)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()
			stats := []statistics{}
			totalDownloads := uint(0)

			err = json.NewDecoder(resp.Body).Decode(&stats)
			if err != nil {
				panic(err)
			}

			for _, stat := range stats {
				if stat.Assets == nil || len(stat.Assets) == 0 {
					continue
				}

				for _, asset := range stat.Assets {
					totalDownloads += asset.DownloadCount
					fmt.Printf("Repo: %s\tAsset: %s\tCount: %d\tDate: %s\n", reponame, asset.Name, asset.DownloadCount, asset.Date.Format(time.RFC850))
				}

				fmt.Println()
			}

			fmt.Printf("Total downloads for %s: %d\n\n", reponame, totalDownloads)
		}(reposApiEndpoint+repo.Name+"/releases", repo.Name)
	}

	wg.Wait()
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 || len(args) > 2 {
		fmt.Println("Usage: gcdstats [user] [repository, optional]")
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
