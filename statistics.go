package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type statistics struct {
	Assets []struct {
		Name          string    `json:"name"`
		DownloadCount uint      `json:"download_count"`
		Date          time.Time `json:"updated_at"`
	} `json:"assets"`
}

func fetchStatistics(repos []repository) {
	wg := &sync.WaitGroup{}
	wg.Add(len(repos))

	const reposApiEndpoint = "https://api.github.com/repos/"
	for _, repo := range repos {
		go fetchStatisticsForRepo(reposApiEndpoint+repo.Name+"/releases", repo.Name, wg)
	}

	wg.Wait()
}

func fetchStatisticsForRepo(repourl, reponame string, wg *sync.WaitGroup) {
	defer wg.Done()

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

	buffer := bytes.Buffer{}

	for _, stat := range stats {
		if stat.Assets == nil || len(stat.Assets) == 0 {
			continue
		}

		for _, asset := range stat.Assets {
			totalDownloads += asset.DownloadCount
			buffer.WriteString(
				fmt.Sprintf("Repo: %s\tAsset: %s\tCount: %d\tDate: %s\n", reponame, asset.Name, asset.DownloadCount, asset.Date.Format(time.RFC850)),
			)
		}

		buffer.WriteRune('\n')
	}

	buffer.WriteString(fmt.Sprintf("Total downloads for %s: %d\n\n", reponame, totalDownloads))

	fmt.Print(buffer.String())
}
