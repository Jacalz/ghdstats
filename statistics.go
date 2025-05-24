package main

import (
	"encoding/json"
	"fmt"
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
		os.Stderr.WriteString("Error: The HTTP get request failed. Error message: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	defer resp.Body.Close()
	stats := []statistics{}
	totalDownloads := uint64(0)

	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		os.Stderr.WriteString("Error: Failed to decode JSON data. A likely culprit is that the GitHub API limit was likely reached.\nTry again in a few hours. Error message: ")
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}

	buffer := make([]byte, 0, 4096)

	for _, stat := range stats {
		if stat.Assets == nil {
			continue
		}

		for _, asset := range stat.Assets {
			totalDownloads += asset.DownloadCount

			buffer = fmt.Appendf(buffer, "Repo: %-20v Asset: %-40v Count: %-5v",
				reponame,
				asset.Name,
				strconv.FormatUint(asset.DownloadCount, 10),
			)

			buffer = asset.Date.AppendFormat(buffer, time.RFC850)
			buffer = append(buffer, '\n')
		}
		buffer = append(buffer, '\n')
	}

	buffer = fmt.Appendf(buffer, "Total downloads for %s: ", reponame)
	buffer = strconv.AppendUint(buffer, totalDownloads, 10)
	buffer = append(buffer, '\n')

	os.Stdout.Write(buffer)
}
