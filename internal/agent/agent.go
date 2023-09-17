package agent

import (
	"alerting/internal/config"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

func CollectMetrics(cfg *config.AgentConfig) {
	var currentStats runtime.MemStats
	var fieldValue string
	counter := 0
	for {
		runtime.ReadMemStats(&currentStats)
		rStats := reflect.ValueOf(currentStats)
		if counter%cfg.ReportInterval == 0 {
			for _, metricName := range cfg.GaugeMetricsList {
				field := rStats.FieldByName(metricName)

				switch field.Interface().(type) {
				case uint64:
					{
						fieldValue = strconv.FormatUint(field.Uint(), 10)
					}
				case float64:
					{
						fieldValue = strconv.FormatFloat(field.Float(), 'f', -1, 64)
					}
				}

				sendMetric(cfg, metricName, fieldValue, "gauge")
			}
			sendMetric(cfg, "PollCount", "5", "counter")
			sendMetric(cfg, "RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64), "gauge")

		}
		counter++
		counter = counter % (cfg.PollInterval * cfg.ReportInterval)
		time.Sleep(1 * time.Second)
	}

}

func sendMetric(cfg *config.AgentConfig, name string, value string, mType string) error {
	sendAddr := "http://" + cfg.ServerURL + "/update/{metricType}/{metricName}/{metricValue}"
	client := resty.New()
	_, err := client.R().SetPathParams(map[string]string{
		"metricName":  name,
		"metricValue": value,
		"metricType":  mType,
	}).
		Post(sendAddr)

	if err != nil {
		return err
	}
	return err

}
