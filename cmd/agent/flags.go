package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	addr           string
	reportInterval int
	pollInterval   int
)

func parseFlags() error {

	flag.StringVar(&addr, "a", "localhost:8080", "Server address")
	flag.IntVar(&reportInterval, "r", 10, "Report interval")
	flag.IntVar(&pollInterval, "p", 2, "Report interval")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		addr = envAddr
	}

	if envReportInt := os.Getenv("REPORT_INTERVAL"); envReportInt != "" {
		envReportIntVal, err := strconv.Atoi(envReportInt)
		if err != nil {
			return err
		}

		reportInterval = envReportIntVal
	}

	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		envPollIntVal, err := strconv.Atoi(envPollInt)
		if err != nil {
			return err
		}

		pollInterval = envPollIntVal
	}

	return nil
}
