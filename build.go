//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type target struct {
	goos   string
	goarch string
}

var targets = []target{
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
	os.MkdirAll("release", 0750)

	for _, t := range targets {
		path := filepath.Join("release", "ghdstats-"+t.goos+"-"+t.goarch)
		if t.goos == "windows" {
			path += ".exe"
		}

		cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w", "-o", path)
		cmd.Env = append(os.Environ(), "GOOS="+t.goos, "GOARCH="+t.goarch)

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s %s err: %s, out: %s\n", t.goos, t.goarch, err, out)
			os.Exit(1)
		}
	}
}
