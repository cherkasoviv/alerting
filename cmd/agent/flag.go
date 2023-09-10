package main

import (
	"flag"
	"os"
	"strconv"
)

var flagRunAddr string
var reportInterval, pollInterval int

func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of reporting metrics to server")
	flag.IntVar(&pollInterval, "p", 2, "frequency of recording metrics in agent")
	flag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	if envReportInterval, err := strconv.Atoi(os.Getenv("REPORT_INTERVAL")); envReportInterval != 0 && err == nil {
		reportInterval = envReportInterval
	}

	if envPollInterval, err := strconv.Atoi(os.Getenv("POLL_INTERVAL")); envPollInterval != 0 && err == nil {
		pollInterval = envPollInterval
	}
}
