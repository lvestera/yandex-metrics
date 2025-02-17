package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Addr           string
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	err := parseFlags(cfg)

	return cfg, err
}

func parseFlags(cfg *Config) error {

	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "Server address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Report interval")
	flag.StringVar(&cfg.Key, "k", "", "Secret Key")
	flag.IntVar(&cfg.RateLimit, "l", 1, "Rate limit")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envReportInt := os.Getenv("REPORT_INTERVAL"); envReportInt != "" {
		envReportIntVal, err := strconv.Atoi(envReportInt)
		if err != nil {
			return err
		}

		cfg.ReportInterval = envReportIntVal
	}

	if envPollInt := os.Getenv("POLL_INTERVAL"); envPollInt != "" {
		envPollIntVal, err := strconv.Atoi(envPollInt)
		if err != nil {
			return err
		}

		cfg.PollInterval = envPollIntVal
	}

	if envKey := os.Getenv("KEY"); envKey != "" {
		cfg.Key = envKey
	}

	if envRateInt := os.Getenv("RATE_LIMIT"); envRateInt != "" {
		envRateIntVal, err := strconv.Atoi(envRateInt)
		if err != nil {
			return err
		}

		cfg.RateLimit = envRateIntVal
	}

	return nil
}
