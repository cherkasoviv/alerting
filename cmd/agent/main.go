package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

const URL = "http://localhost:8080/update/"

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

func main() {
	var currentStats runtime.MemStats
	var fieldValue string
	counter := 0
	for {
		runtime.ReadMemStats(&currentStats)
		rStats := reflect.ValueOf(currentStats)
		if counter%5 == 0 {
			for _, metricName := range RuntimeGaugeMetrics {
				field := rStats.FieldByName(metricName)

				switch field.Interface().(type) {
				case uint64:
					{
						fieldValue = strconv.FormatUint(field.Uint(), 10)
					}
				case float64:
					{
						fieldValue = strconv.FormatFloat(field.Float(), 'f', 2, 64)
					}
				}

				sendMetric(metricName, fieldValue, "gauge")
			}
			sendMetric("PollCount", "5", "counter")
			sendMetric("RandomValue", strconv.FormatFloat(rand.Float64(), 'f', 2, 64), "gauge")

		}
		counter++
		time.Sleep(2 * time.Second)
	}
}

func sendMetric(name string, value string, mType string) error {
	url := URL + mType + "/" + name + "/" + value
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp)
	defer resp.Body.Close()

	if err != nil {
		return err
	}
	return nil
}
