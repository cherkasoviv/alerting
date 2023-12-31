package metrics

import "strconv"

type GaugeMetric struct {
	Metric Metric
	Value  float64
}

func (gaugeMetric *GaugeMetric) GetName() string {
	return gaugeMetric.Metric.Name
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

func (gaugeMetric *GaugeMetric) String() string {
	return gaugeMetric.GetName() + ":" + strconv.FormatFloat(gaugeMetric.Value, 'f', -1, 64)
}

func (gaugeMetric *GaugeMetric) GetValue() string {
	return strconv.FormatFloat(gaugeMetric.Value, 'f', -1, 64)
}

func (gaugeMetric *GaugeMetric) GetType() string {
	return string(gaugeMetric.Metric.Mtype)
}
