package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounterMetric_GetName(t *testing.T) {
	type fields struct {
		CMetric Metric
		value   uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple_test",
			fields: struct {
				CMetric Metric
				value   uint64
			}{CMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "counter", Mtype: Counter}, value: 0},
			want: "counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterMetric := &CounterMetric{
				CMetric: tt.fields.CMetric,
				value:   tt.fields.value,
			}

			assert.Equal(t, tt.want, counterMetric.GetName())

		})
	}
}

func TestCounterMetric_GetValue(t *testing.T) {
	type fields struct {
		CMetric Metric
		value   uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple_test",
			fields: struct {
				CMetric Metric
				value   uint64
			}{CMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "counter", Mtype: Counter}, value: 0},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterMetric := &CounterMetric{
				CMetric: tt.fields.CMetric,
				value:   tt.fields.value,
			}
			assert.Equal(t, tt.want, counterMetric.GetValue())
		})
	}
}

func TestCounterMetric_String(t *testing.T) {
	type fields struct {
		CMetric Metric
		value   uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple_test",
			fields: struct {
				CMetric Metric
				value   uint64
			}{CMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "counter", Mtype: Counter}, value: 0},
			want: "counter:0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterMetric := &CounterMetric{
				CMetric: tt.fields.CMetric,
				value:   tt.fields.value,
			}
			assert.Equal(t, tt.want, counterMetric.String())
		})
	}
}

func TestCounterMetric_UpdateValue(t *testing.T) {
	type fields struct {
		CMetric Metric
		value   uint64
	}
	type args struct {
		newValue string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		newFields fields
	}{
		{
			name: "simple_test",
			fields: struct {
				CMetric Metric
				value   uint64
			}{CMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "counter", Mtype: Counter}, value: 1},
			newFields: struct {
				CMetric Metric
				value   uint64
			}{CMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "counter", Mtype: Counter}, value: 6},
			args: struct{ newValue string }{newValue: "5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterMetric := &CounterMetric{
				CMetric: tt.fields.CMetric,
				value:   tt.fields.value,
			}
			updatedMetric := &CounterMetric{
				CMetric: tt.newFields.CMetric,
				value:   tt.newFields.value,
			}
			counterMetric.UpdateValue(tt.args.newValue)
			assert.Equal(t, updatedMetric, counterMetric)
		})
	}
}
