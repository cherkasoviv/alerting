package main

import (
	"alerting/internal/agent"
	"alerting/internal/config"
)

func main() {

	cfg := config.LoadAgentConfig()
	jobs := make(chan agent.MetricJob, 1024)
	for i := 0; i < cfg.RateLimit; i++ {
		go agent.SendMetricJSON(cfg, jobs)
	}
	go agent.CollectGopsUtilMetrics(cfg, jobs)
	agent.CollectMetrics(cfg, jobs)

}
