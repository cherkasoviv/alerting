package main

import (
	"alerting/internal/config"
	"alerting/internal/utils"
)

func main() {

	cfg := config.LoadAgentConfig()
	utils.CollectMetrics(cfg)

}
