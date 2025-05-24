//go:build profile

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go log.Println(http.ListenAndServe("localhost:6060", nil))
}
