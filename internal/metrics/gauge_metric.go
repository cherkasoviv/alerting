package metrics

import "strconv"

type GaugeMetric struct {
	GMetric Metric
	Value   float64
}

func (gaugeMetric *GaugeMetric) GetName() string {
	return gaugeMetric.GMetric.Name
}

func (gaugeMetric *GaugeMetric) UpdateValue(newValue string) error {

	newFloatValue, err := strconv.ParseFloat(newValue, 64)

	if err == nil {
		gaugeMetric.Value = newFloatValue
		return nil
	} else {
		return err
	}
}

func (gaugeMetric *GaugeMetric) String() (string, error) {
	return gaugeMetric.GetName() + ":" + strconv.FormatFloat(gaugeMetric.Value, 'f', 2, 64), nil
}
