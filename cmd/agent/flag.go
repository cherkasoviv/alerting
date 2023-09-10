package main

import "flag"

var flagRunAddr string
var reportInterval, pollInterval int

func parseFlags() {

	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of reporting metrics to server")
	flag.IntVar(&pollInterval, "p", 2, "frequency of recording metrics in agent")
	flag.Parse()
}
