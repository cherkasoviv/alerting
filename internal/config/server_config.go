package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Host string
}

func (cfg *ServerConfig) parseFlags() {

	flag.StringVar(&cfg.Host, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

}

func (cfg *ServerConfig) parseEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.Host = envRunAddr
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
