package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type statistics struct {
	Assets []struct {
		Name          string    `json:"name"`
		DownloadCount uint64    `json:"download_count"`
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
	totalDownloads := uint64(0)

	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		panic(err)
	}

	buffer := make([]byte, 0, 30+len(reponame))

	for _, stat := range stats {
		if stat.Assets == nil || len(stat.Assets) == 0 {
			continue
		}

		for _, asset := range stat.Assets {
			totalDownloads += asset.DownloadCount

			buffer = append(buffer, "Repo: "...)
			buffer = append(buffer, reponame...)
			buffer = append(buffer, '\t')

			buffer = append(buffer, "Asset: "...)
			buffer = append(buffer, asset.Name...)
			buffer = append(buffer, '\t', '\t')

			buffer = append(buffer, "Count: "...)
			buffer = strconv.AppendUint(buffer, asset.DownloadCount, 10)
			buffer = append(buffer, '\t')

			buffer = append(buffer, "Date: "...)
			buffer = asset.Date.AppendFormat(buffer, time.RFC850)
			buffer = append(buffer, '\n')
		}

		buffer = append(buffer, '\n')
	}

	buffer = append(buffer, "Total downloads for "...)
	buffer = append(buffer, reponame...)
	buffer = append(buffer, ':', ' ')
	buffer = strconv.AppendUint(buffer, totalDownloads, 10)
	buffer = append(buffer, '\n', '\n')

	os.Stdout.Write(buffer)
}
