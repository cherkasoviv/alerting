package metrics

import (
	"strconv"
)

type CounterMetric struct {
	CMetric Metric
	value   uint64
}

func (counterMetric *CounterMetric) GetName() string {
	return counterMetric.CMetric.Name
}
func (counterMetric *CounterMetric) UpdateValue(newValue string) error {
	newIntValue, err := strconv.ParseUint(newValue, 10, 64)

	if err != nil {
		return err
	}
	counterMetric.value += newIntValue
	return nil
}

func (counterMetric *CounterMetric) String() string {
	return counterMetric.GetName() + ":" + strconv.FormatUint(counterMetric.value, 10)
}

func (counterMetric *CounterMetric) GetValue() string {
	return strconv.FormatUint(counterMetric.value, 10)
}
