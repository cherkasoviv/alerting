package main

import (
	"github.com/go-resty/resty/v2"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
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

	client := resty.New()
	_, err := client.R().SetPathParams(map[string]string{
		"metricName":  name,
		"metricValue": value,
		"metricType":  mType,
	}).
		Post("http://localhost:8080/update/{metricType}/{metricName}/{metricValue}")

	if err != nil {
		return err
	}
	return err

}
