//go:build !profile

package main

func profile() func() {
	return nil
}
