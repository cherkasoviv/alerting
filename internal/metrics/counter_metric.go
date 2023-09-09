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

	if err == nil {
		counterMetric.value += newIntValue
		return nil
	} else {
		return err
	}
}

func (counterMetric CounterMetric) String() (string, error) {
	return counterMetric.GetName() + ":" + strconv.FormatUint(counterMetric.value, 10), nil
}
