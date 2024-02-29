package main

import (
	"alerting/internal/agent"
	"alerting/internal/config"
)

func main() {

	cfg := config.LoadAgentConfig()
	agent.CollectMetrics(cfg)

}
