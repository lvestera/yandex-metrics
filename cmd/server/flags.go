package main

import (
	"flag"
	"os"
)

var (
	addr string
)

func parseFlags() {

	flag.StringVar(&addr, "a", "localhost:8081", "Server address")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		addr = envAddr
	}
}
