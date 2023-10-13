package metrics

import (
	"strconv"
)

type CounterMetric struct {
	Metric Metric
	Value  uint64
}

func (counterMetric *CounterMetric) GetName() string {
	return counterMetric.Metric.Name
}
func (counterMetric *CounterMetric) UpdateValue(newValue string) error {
	newIntValue, err := strconv.ParseUint(newValue, 10, 64)

	if err != nil {
		return err
	}
	counterMetric.Value += newIntValue
	return nil
}

func (counterMetric *CounterMetric) String() string {
	return counterMetric.GetName() + ":" + strconv.FormatUint(counterMetric.Value, 10)
}

func (counterMetric *CounterMetric) GetValue() string {
	return strconv.FormatUint(counterMetric.Value, 10)
}

func (counterMetric *CounterMetric) GetType() string {
	return string(counterMetric.Metric.Mtype)
}
