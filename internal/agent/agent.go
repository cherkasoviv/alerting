package agent

import (
	"alerting/internal/config"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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

				sendMetricJSON(cfg, metricName, fieldValue, "gauge")
			}

			sendMetricJSON(cfg, "PollCount", "5", "counter")

			sendMetricJSON(cfg, "RandomValue", strconv.FormatFloat(rand.Float64(), 'f', -1, 64), "gauge")
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

func sendMetricJSON(cfg *config.AgentConfig, name string, value string, mType string) error {
	sendAddr := "http://" + cfg.ServerURL + "/update/"
	client := resty.New()
	metricToSend := metrics{
		ID:    name,
		MType: mType,
		Delta: nil,
		Value: nil,
	}

	switch mType {
	case "gauge":
		{
			floatValue, _ := strconv.ParseFloat(value, 64)
			metricToSend.Value = &floatValue
		}
	case "counter":
		{
			intValue, _ := strconv.ParseInt(value, 10, 64)
			metricToSend.Delta = &intValue
		}
	}
	req, _ := json.Marshal(metricToSend)
	_, err := client.R().SetBody(req).
		Post(sendAddr)

	if err != nil {
		return err
	}
	return err

}
