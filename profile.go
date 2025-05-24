//go:build profile

package main

import (
	"log"
	"os"
	"runtime/pprof"
)

func profile() func() {
	f, err := os.Create("default.pgo")
	if err != nil {
		log.Fatal("Could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("Could not start CPU profile: ", err)
	}
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}
