package metrics

import "strconv"

type CounterMetric struct {
	CMetric Metric
	value   int64
}

func (counterMetric *CounterMetric) GetName() string {
	return counterMetric.CMetric.Name
}
func (counterMetric *CounterMetric) UpdateValue(newValue string) error {

	newIntValue, err := strconv.ParseInt(newValue, 10, 64)

	if err == nil {
		counterMetric.value = newIntValue
		return nil
	} else {
		return err
	}
}
