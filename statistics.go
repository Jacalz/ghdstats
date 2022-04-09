package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"os"
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

	buffer := strings.Builder{}

	for _, stat := range stats {
		if stat.Assets == nil || len(stat.Assets) == 0 {
			continue
		}

		for _, asset := range stat.Assets {
			totalDownloads += asset.DownloadCount

			buffer.WriteString("Repo: ")
			buffer.WriteString(reponame)
			buffer.WriteByte('\t')

			buffer.WriteString("Asset: ")
			buffer.WriteString(asset.Name)
			buffer.WriteByte('\t')

			buffer.WriteString("Count: ")
			buffer.WriteString(strconv.FormatUint(asset.DownloadCount, 10))
			buffer.WriteByte('\t')

			buffer.WriteString("Date: ")
			buffer.WriteString(asset.Date.Format(time.RFC850))
			buffer.WriteByte('\n')
		}

		buffer.WriteByte('\n')
	}

	buffer.WriteString("Total downloads for ")
	buffer.WriteString(reponame)
	buffer.WriteString(": ")
	buffer.WriteString(strconv.FormatUint(totalDownloads, 10))
	buffer.WriteString("\n\n")

	os.Stdout.WriteString(buffer.String())
}
