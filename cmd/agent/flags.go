package main

import (
	"flag"
)

var (
	addr           string
	reportInterval int
	pollInterval   int
)

func parseFlags() {

	flag.StringVar(&addr, "a", "localhost:8080", "Server address")
	flag.IntVar(&reportInterval, "r", 10, "Report interval")
	flag.IntVar(&pollInterval, "p", 2, "Report interval")
	flag.Parse()
}
