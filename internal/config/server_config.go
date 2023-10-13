package config

import (
	"flag"
	"os"
	"strconv"
)

type ServerConfig struct {
	Host            string
	StoreInterval   int
	FileStoragePath string
	NeedToRestore   bool
}

func (cfg *ServerConfig) parseFlags() {

	flag.StringVar(&cfg.Host, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.StoreInterval, "i", 300, "backup period for inmemory db")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metrics-db.json", "backup path")
	flag.BoolVar(&cfg.NeedToRestore, "r", true, "need backup true/false")
	flag.Parse()

	if cfg.FileStoragePath != "" {
		cfg.FileStoragePath = "." + cfg.FileStoragePath
	}

}

func (cfg *ServerConfig) parseEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Host = envRunAddr
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		interval, err := strconv.ParseInt(envStoreInterval, 10, 0)
		if err == nil {
			cfg.StoreInterval = int(interval)
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = "." + envFileStoragePath
	}

	if envNeedToRestore := os.Getenv("RESTORE"); envNeedToRestore != "" {
		restore, err := strconv.ParseBool(envNeedToRestore)
		if err == nil {
			cfg.NeedToRestore = restore
		}
	}
}

func LoadServerConfig() *ServerConfig {
	cfg := ServerConfig{
		Host: "",
	}
	cfg.parseFlags()
	cfg.parseEnv()
	return &cfg

}
