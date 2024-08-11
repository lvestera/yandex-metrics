package main

import (
	"flag"
)

var (
	addr string
)

func parseFlags() {

	flag.StringVar(&addr, "a", "localhost:8080", "Server address")
	flag.Parse()
}
