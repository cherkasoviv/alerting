package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGaugeMetric_GetName(t *testing.T) {
	type fields struct {
		GMetric Metric
		value   float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple_test",
			fields: struct {
				GMetric Metric
				value   float64
			}{GMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "gauge", Mtype: Gauge}, value: 0},
			want: "gauge",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gaugeMetric := &GaugeMetric{
				Metric: tt.fields.GMetric,
				Value:  tt.fields.value,
			}

			assert.Equal(t, tt.want, gaugeMetric.GetName())

		})
	}
}

func TestGaugeMetric_GetValue(t *testing.T) {
	type fields struct {
		GMetric Metric
		value   float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple_test",
			fields: struct {
				GMetric Metric
				value   float64
			}{GMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "gauge", Mtype: Gauge}, value: 0},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gaugeMetric := &GaugeMetric{
				Metric: tt.fields.GMetric,
				Value:  tt.fields.value,
			}
			assert.Equal(t, tt.want, gaugeMetric.GetValue())
		})
	}
}

func TestGaugeMetric_String(t *testing.T) {
	type fields struct {
		GMetric Metric
		value   float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple_test",
			fields: struct {
				GMetric Metric
				value   float64
			}{GMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "gauge", Mtype: Gauge}, value: 0},
			want: "gauge:0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gaugeMetric := &GaugeMetric{
				Metric: tt.fields.GMetric,
				Value:  tt.fields.value,
			}
			assert.Equal(t, tt.want, gaugeMetric.String())
		})
	}
}

func TestGaugeMetric_UpdateValue(t *testing.T) {
	type fields struct {
		GMetric Metric
		value   float64
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
				GMetric Metric
				value   float64
			}{GMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "gauge", Mtype: Gauge}, value: 0},
			newFields: struct {
				GMetric Metric
				value   float64
			}{GMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "gauge", Mtype: Gauge}, value: 5},
			args: struct{ newValue string }{newValue: "5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gaugeMetric := &GaugeMetric{
				Metric: tt.fields.GMetric,
				Value:  tt.fields.value,
			}
			updatedMetric := &GaugeMetric{
				Metric: tt.newFields.GMetric,
				Value:  tt.newFields.value,
			}
			gaugeMetric.UpdateValue(tt.args.newValue)
			assert.Equal(t, updatedMetric, gaugeMetric)
		})
	}
}
