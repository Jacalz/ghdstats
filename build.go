//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

const version = "v1.2.0"

type target struct {
	os   string
	arch string
}

var targets = [...]target{
	{"linux", "amd64"},
	{"linux", "arm64"},

	{"darwin", "amd64"},
	{"darwin", "arm64"},

	{"windows", "amd64"},
	{"windows", "arm64"},

	{"freebsd", "amd64"},
	{"freebsd", "arm64"},
	{"openbsd", "amd64"},
	{"openbsd", "arm64"},
	{"netbsd", "amd64"},
	{"netbsd", "arm64"},
}

func main() {
	err := os.MkdirAll("builds", 0o750)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create builds directory: %v\n", err)
		os.Exit(1)
	}

	wg := errgroup.Group{}
	for _, tar := range targets {
		wg.Go(func() error {
			path := filepath.Join("builds", "ghdstats-"+version+"-"+tar.os+"-"+tar.arch)
			if tar.os == "windows" {
				path += ".exe"
			}

			cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", path)
			cmd.Env = append(os.Environ(), "GOOS="+tar.os, "GOARCH="+tar.arch)
			if tar.arch == "amd64" {
				cmd.Env = append(cmd.Env, "GOAMD64=v3")
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		})
	}

	if err := wg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Build failed: %v\n", err)
		os.Exit(1)
	}
}
