package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Addr            string
	StorageInterval int
	FileStoragePath string
	Restore         bool
	DBConfig        string
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	err := parseFlags(cfg)

	return cfg, err
}

func parseFlags(cfg *Config) error {

	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "Server address")
	flag.IntVar(&cfg.StorageInterval, "i", 300, "Storage interval")
	flag.StringVar(&cfg.FileStoragePath, "f", "file.txt", "File storage path")
	flag.BoolVar(&cfg.Restore, "r", true, "Restore data on server run")
	flag.StringVar(&cfg.DBConfig, "d", "", "Database connection")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}

	if envStorageInterval := os.Getenv("STORE_INTERVAL"); envStorageInterval != "" {
		envStorageIntervalVal, err := strconv.Atoi(envStorageInterval)
		if err != nil {
			return err
		}

		cfg.StorageInterval = envStorageIntervalVal
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		envRestoreVal, err := strconv.ParseBool(envRestore)
		if err != nil {
			return err
		}

		cfg.Restore = envRestoreVal
	}

	if envDBConfig := os.Getenv("DATABASE_DSN"); envDBConfig != "" {
		cfg.DBConfig = envDBConfig
	}

	return nil
}
