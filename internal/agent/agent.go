package agent

import (
	"alerting/internal/config"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	cpu2 "github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"
)

type MetricJob struct {
	ID    string
	MType string
	Value string
}
type metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func CollectMetrics(cfg *config.AgentConfig, jobChannel chan<- MetricJob) {
	var currentStats runtime.MemStats
	counter := 0
	for {
		runtime.ReadMemStats(&currentStats)
		rStats := reflect.ValueOf(currentStats)
		if counter%cfg.ReportInterval == 0 {
			for _, metricName := range cfg.GaugeMetricsList {
				field := rStats.FieldByName(metricName)
				var metric MetricJob
				metric.ID = metricName
				metric.MType = "gauge"
				switch field.Interface().(type) {
				case uint64:
					{
						metric.Value = strconv.FormatUint(field.Uint(), 10)
					}
				case float64:
					{
						metric.Value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
					}
				}
				jobChannel <- metric

			}
			polMetricJob := MetricJob{
				ID:    "PollCount",
				MType: "counter",
				Value: "5",
			}
			jobChannel <- polMetricJob

			randomMetricJob := MetricJob{
				ID:    "RandomValue",
				MType: "gauge",
				Value: strconv.FormatFloat(rand.Float64(), 'f', -1, 64),
			}
			jobChannel <- randomMetricJob
		}
		counter++
		counter = counter % (cfg.PollInterval * cfg.ReportInterval)
		time.Sleep(1 * time.Second)
	}

}
func CollectGopsUtilMetrics(cfg *config.AgentConfig, jobChannel chan<- MetricJob) {
	counter := 0
	for {

		if counter%cfg.ReportInterval == 0 {
			vm, _ := mem.VirtualMemory()
			TotalMemoryMetricJob := MetricJob{
				ID:    "TotalMemory",
				MType: "gauge",
				Value: strconv.Itoa(int(vm.Total)),
			}
			jobChannel <- TotalMemoryMetricJob

			FreeMemoryMetricJob := MetricJob{
				ID:    "FreeMemory",
				MType: "gauge",
				Value: strconv.Itoa(int(vm.Free)),
			}
			jobChannel <- FreeMemoryMetricJob

			cpu, _ := cpu2.Percent(0, false)
			CPUutilization := MetricJob{
				ID:    "CPUutilization",
				MType: "gauge",
				Value: strconv.FormatFloat(cpu[0], 'f', 2, 64),
			}
			jobChannel <- CPUutilization

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
	}).SetHeader("Accept-Encoding", "gzip").
		Post(sendAddr)

	if err != nil {
		return err
	}
	return err

}

func SendMetricJSON(cfg *config.AgentConfig, jobs <-chan MetricJob) {
	sendAddr := "http://" + cfg.ServerURL + "/update/"
	client := resty.New()
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case mJob := <-jobs:
			{
				metricToSend := metrics{
					ID:    mJob.ID,
					MType: mJob.MType,
					Delta: nil,
					Value: nil,
				}

				switch mJob.MType {
				case "gauge":
					{
						floatValue, _ := strconv.ParseFloat(mJob.Value, 64)
						metricToSend.Value = &floatValue
					}
				case "counter":
					{
						intValue, _ := strconv.ParseInt(mJob.Value, 10, 64)
						metricToSend.Delta = &intValue
					}
				}
				req, _ := json.Marshal(metricToSend)
				_, err := client.R().SetBody(req).
					Post(sendAddr)

				if err != nil {
					panic(err)
				}
			}
		case <-ticker.C:
			{
				continue
			}
		}
	}

}
