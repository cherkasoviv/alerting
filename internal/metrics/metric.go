package metrics

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type (
	AbstractMetric interface {
		UpdateValue(newValue string) error
		GetName() string
		String() string
		GetValue() string
	}
)

type Metric struct {
	Name  string
	Mtype MetricType
}
