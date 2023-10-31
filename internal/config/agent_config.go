package config

import (
	"flag"
	"os"
	"strconv"
)

var RuntimeGaugeMetrics = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

type AgentConfig struct {
	ServerURL        string
	ReportInterval   int
	PollInterval     int
	GaugeMetricsList []string
	HashSHA256Key    string
	RateLimit        int
}

func (cfg *AgentConfig) parseFlags() {

	flag.StringVar(&cfg.ServerURL, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "frequency of reporting metrics to server")
	flag.IntVar(&cfg.PollInterval, "p", 2, "frequency of recording metrics in agent")
	flag.StringVar(&cfg.HashSHA256Key, "k", "", "hash key")
	flag.IntVar(&cfg.RateLimit, "l", 1, "number of metric senders")
	flag.Parse()

}

func (cfg *AgentConfig) parseEnv() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.ServerURL = envRunAddr
	}

	if envReportInterval, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL")); envReportInterval != 0 && err == nil {
		cfg.ReportInterval = envReportInterval
	}

	if envPollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL")); envPollInterval != 0 && err == nil {
		cfg.PollInterval = envPollInterval
	}

	if envHashSHA256Key := os.Getenv("KEY"); envHashSHA256Key != "" {
		cfg.HashSHA256Key = envHashSHA256Key
	}

	if envRateLimit, err := strconv.Atoi(os.Getenv("RATE_LIMIT")); envRateLimit != 0 && err == nil {
		cfg.RateLimit = envRateLimit
	}
}

func LoadAgentConfig() *AgentConfig {
	cfg := AgentConfig{
		ServerURL:        "",
		ReportInterval:   0,
		PollInterval:     0,
		GaugeMetricsList: []string{},
		HashSHA256Key:    "",
		RateLimit:        1,
	}

	cfg.parseFlags()
	cfg.parseEnv()
	cfg.GaugeMetricsList = RuntimeGaugeMetrics

	return &cfg
}
