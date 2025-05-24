package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

const reposApiEndpoint = "https://api.github.com/repos/"

type statistics struct {
	Assets []struct {
		Name          string    `json:"name"`
		DownloadCount uint64    `json:"download_count"`
		Date          time.Time `json:"updated_at"`
	} `json:"assets"`
}

func fetchStatistics(repos []repository) error {
	if len(repos) == 1 {
		return fetchStatisticsForRepo(repos[0].Name)
	}

	wg := errgroup.Group{}
	for _, repo := range repos {
		wg.Go(func() error {
			return fetchStatisticsForRepo(repo.Name)
		})
	}

	return wg.Wait()
}

func fetchStatisticsForRepo(reponame string) error {
	resp, err := http.Get(reposApiEndpoint + reponame + "/releases")
	if err != nil {
		return err
	} else if resp.StatusCode == http.StatusForbidden {
		return errorRateLimitExceeded
	}

	defer resp.Body.Close()
	stats := []statistics{}
	totalDownloads := uint64(0)

	err = json.NewDecoder(resp.Body).Decode(&stats)
	if err != nil {
		return err
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

	_, err = os.Stdout.Write(buffer)
	return err
}
